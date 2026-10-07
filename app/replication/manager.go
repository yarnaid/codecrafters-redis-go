package replication

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"strconv"
	"strings"
	"sync"

	"my-redis/app/parser"
)

const replicaQueueSize = 1024

type Replica struct {
	conn   net.Conn
	out    chan []byte
	closed chan struct{}
	once   sync.Once
	logger *slog.Logger
	ack    int
}

func (r *Replica) Close() {
	r.once.Do(func() {
		close(r.closed)
		if err := r.conn.Close(); err != nil {
			slog.Error("problem with closing replica conn", "err", err, "addr", r.conn.RemoteAddr())
		}
	})
}

type Manager struct {
	mu       sync.Mutex
	replicas map[*Replica]struct{}
	offset   int
	replID   ReplID
	role     Role
	logger   *slog.Logger
}

func (m *Manager) ReplID() ReplID {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.replID
}

// Register is called from the PSYNC command after FULLRESYNC is decided.
// Snapshot and registration happen under the same lock that Propagate uses,
// so no write is lost between the RDB and the first streamed command.
func (m *Manager) Register(conn net.Conn, snapshot func() []byte, reader *parser.Reader) *Replica {
	m.logger.Info("registering new replica", "addr", conn.RemoteAddr())
	r := &Replica{conn: conn, out: make(chan []byte, replicaQueueSize), closed: make(chan struct{}), logger: m.logger.With("addr", conn.RemoteAddr())}

	m.mu.Lock()
	// m.logger.Debug("preparing full resync", "addr", conn.RemoteAddr())
	out, _ := parser.SimpleString(fmt.Sprintf("FULLRESYNC %s 0", m.replID)).Serialize()
	out = append(out, snapshot()...)
	r.out <- out // first item in the queue, before any propagated command
	// m.logger.Debug("snapshot sent to replica", "addr", conn.RemoteAddr())
	m.replicas[r] = struct{}{}
	m.mu.Unlock()

	go m.writeLoop(r)
	go m.readLoop(r, reader)
	// go func() {
	// 	timer := time.After(time.Second * 1)
	// 	<-timer
	// 	m.logger.Debug("sending PING to replicas")
	// 	data, _ := parser.CommandFromStrings("PING").Serialize()
	// 	r.conn.Write(data)
	// 	m.logger.Debug("sending ACK to replicas")
	// 	data, _ = parser.CommandFromStrings("REPLCONF", "GETACK", "*").Serialize()
	// 	r.conn.Write(data)
	// }()
	return r
}

// AddProcessed advances the replica offset by n bytes of the master stream.
func (m *Manager) AddProcessed(n int) {
	m.mu.Lock()
	m.offset += n
	m.mu.Unlock()
}

// AdoptMaster records the master's identity after FULLRESYNC (replica side).
func (m *Manager) AdoptMaster(replID ReplID, offset int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.role = ReplicaRole
	m.replID = replID
	m.offset = offset
}

func (m *Manager) Offset() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.offset
}

func (m *Manager) Role() Role {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.role
}

func (m *Manager) remove(r *Replica) {
	r.logger.Debug("removing client")
	m.mu.Lock()
	delete(m.replicas, r)
	m.mu.Unlock()
	r.Close()
}

func (m *Manager) writeLoop(r *Replica) {
	r.logger.Info("start replica write loop")
	defer m.remove(r)
	for {
		select {
		case b := <-r.out:
			// m.logger.Debug("writing to replica", "addr", r.conn.RemoteAddr())
			if len(b) == 0 {
				r.logger.Warn("sending zero bytes")
			}
			if _, err := r.conn.Write(b); err != nil {
				r.logger.Error("cannot write to replica", "err", err)
				return
			}
		case <-r.closed:
			r.logger.Debug("replica connection closed")
			return
		}
	}
}

// readLoop consumes REPLCONF ACK <offset> frames from the replica.
func (m *Manager) readLoop(r *Replica, reader *parser.Reader) {
	defer m.remove(r)
	for {
		arr, err := reader.ReadArrays()
		if err != nil {
			if !errors.Is(err, io.EOF) {
				r.logger.Error("cannot read cmds", "err", err)
			}
			r.logger.Info("closing replica connection")
			return
		}
		r.logger.Debug("got from replica", "arr", arr)
		if len(arr) == 0 {
			r.logger.Error("received empty command from replica")
			return
		}
		for _, args := range arr {
			if len(args) == 3 && strings.EqualFold(args[0], "REPLCONF") && strings.EqualFold(args[1], "ACK") {
				if n, err := strconv.Atoi(args[2]); err == nil {
					r.logger.Debug("received ACK", "ACK", n)
					r.ack = n
				}
			}
		}
	}
}

func (m *Manager) Propagate(raw []byte) {
	if len(raw) == 0 {
		m.logger.Warn("attempt to propagate empty cmd")
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.offset += int(len(raw))
	for r := range m.replicas {
		select {
		case r.out <- raw: // shared read-only slice; do not mutate
		default:
			delete(m.replicas, r)
			r.Close() // queue full: replica is too slow, drop it and let it resync
		}
	}
}

func (m *Manager) Count() int {
	return len(m.replicas)
}

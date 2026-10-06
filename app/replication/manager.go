package replication

import (
	"fmt"
	"log/slog"
	"net"
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
	replID   replID
	role     Role
	logger   *slog.Logger
}

func (m *Manager) ReplID() replID {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.replID
}

// Register is called from the PSYNC command after FULLRESYNC is decided.
// Snapshot and registration happen under the same lock that Propagate uses,
// so no write is lost between the RDB and the first streamed command.
func (m *Manager) Register(conn net.Conn, snapshot func() []byte) *Replica {
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
	// go m.readLoop(r)
	return r
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
	delete(m.replicas, r)
}

func (m *Manager) writeLoop(r *Replica) {
	r.logger.Info("start replica write loop")
	defer m.remove(r)
	for {
		select {
		case b := <-r.out:
			// m.logger.Debug("writing to replica", "addr", r.conn.RemoteAddr())
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
// func (m *Manager) readLoop(r *Replica, br *bufio.Reader) {
// 	defer m.remove(r)
// 	for {
// 		args, _, err := parser.ReadCommand(br)
// 		if err != nil {
// 			return
// 		}
// 		if len(args) == 3 && strings.EqualFold(args[0], "REPLCONF") && strings.EqualFold(args[1], "ACK") {
// 			if n, err := strconv.ParseInt(args[2], 10, 64); err == nil {
// 				r.ack.Store(n)
// 			}
// 		}
// 	}
// }

func (m *Manager) Propagate(raw []byte) {
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

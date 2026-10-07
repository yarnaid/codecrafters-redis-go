package server

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/netip"
	"os"
	"strconv"
	"strings"
	"time"

	"my-redis/app/commands"
	"my-redis/app/config"
	"my-redis/app/handler"
	"my-redis/app/parser"
	"my-redis/app/replication"
)

type server struct {
	bind netip.Addr
	port config.Port
	env  *commands.Env
}

func (s *server) listenAddr() netip.AddrPort {
	return netip.AddrPortFrom(s.bind, uint16(s.port))
}

func (s *server) Serve(ctx context.Context) {
	l, err := net.Listen("tcp", s.listenAddr().String())
	if err != nil {
		logger.Error("Failed to bind", "err", err.Error(), "addr", s.bind, "port", s.port)
		os.Exit(1)
	}
	go func() {
		if err := s.listenMaster(ctx); err != nil {
			logger.Error("stop listening master", "err", err.Error())
		}
	}()
	logger.Debug("[server] start listening", "addr", s.bind.String(), "port", s.port, "replica", s.env.Repl.Role() != replication.ReplicaRole, "id", s.env.Repl.ReplID())
	defer context.AfterFunc(ctx, func() {
		if err := l.Close(); err != nil {
			slog.Info("Spot listening port with error", "port", s.port, "err", err)
		}
	})()

	for {
		conn, err := l.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				logger.Info("Connection closed")
				os.Exit(0)
			}
			logger.Error("Error accepting connection", "error", err.Error())
			os.Exit(1)
		}

		h := handler.New(s.env, false)
		go h.HandleConnection(ctx, conn, nil)
	}
}

func (s *server) listenMaster(ctx context.Context) error {
	if s.env.Repl.Role() == replication.MasterRole {
		return nil
	}
	slog.Debug("[server][listenMaster] start connecting")
	conn, err := s.dialMaster(ctx)
	if err != nil {
		return fmt.Errorf("cannot dial master: %w", err)
	}

	line, br, err := s.handshake(conn)
	if err != nil {
		if err1 := conn.Close(); err1 != nil {
			logger.Error("cannot init master conn", "err", err, "err2", err1)
		}
		return fmt.Errorf("cannot handshake master: %w", err)
	}
	_ = conn.SetDeadline(time.Time{})
	err = s.adoptMaster(line)
	if err != nil {
		return fmt.Errorf("cannot adopt master: %w", err)
	}

	slog.Debug("[listenMaster] start replication from master")

	h := handler.New(s.env, true)
	h.HandleConnection(ctx, conn, br)
	return fmt.Errorf("master connection handler stopped")
}

func (s *server) dialMaster(ctx context.Context) (net.Conn, error) {
	dialCtx, cancel := context.WithTimeoutCause(ctx, 5*time.Second, fmt.Errorf("master connection timeout"))
	defer cancel()

	masterAddr, err := s.env.MasterAddr()
	if err != nil {
		return nil, fmt.Errorf("cannot parse master address: %w", err)
	}
	var d net.Dialer
	conn, err := d.DialContext(dialCtx, "tcp", masterAddr)
	if err != nil {
		return nil, fmt.Errorf("dial master: %w", err)
	}
	return conn, nil
}

func (s *server) adoptMaster(line string) error {
	fields := strings.Fields(line)
	if len(fields) != 3 || fields[0][1:] != "FULLRESYNC" {
		return fmt.Errorf("incorrect response from PSYNC: %s", line)
	}
	offset, err := strconv.Atoi(fields[2])
	if err != nil {
		return err
	}
	s.env.Repl.AdoptMaster(replication.ReplID(fields[1]), offset)
	return nil
}

func (s *server) handshake(conn net.Conn) (string, *bufio.Reader, error) {
	br := bufio.NewReader(conn)
	if _, err := sendToMaster(conn, br, "PING"); err != nil {
		return "", nil, err
	}
	if _, err := sendToMaster(conn, br, "REPLCONF", "listening-port", fmt.Sprint(s.port)); err != nil {
		return "", nil, err
	}
	if _, err := sendToMaster(conn, br, "REPLCONF", "capa", "psync2"); err != nil {
		return "", nil, err
	}
	var err error
	var line string
	if line, err = sendToMaster(conn, br, "PSYNC", "?", "-1"); err != nil {
		return "", nil, err
	}
	hdr, err := br.ReadString('\n')
	if err != nil || hdr[0] != '$' {
		return "", nil, fmt.Errorf("handshake RDB header: %w", err)
	}
	n, err := strconv.ParseInt(strings.TrimSpace(hdr[1:]), 10, 64)
	if err != nil || n < 0 {
		return "", nil, errors.New("handshake: bad RDB length")
	}
	_, err = io.CopyN(io.Discard, br, n)
	s.env.Log.Debug("synced to master")
	return line, br, err
}

func sendToMaster(conn net.Conn, br *bufio.Reader, cmd ...string) (string, error) {
	cmdArr := parser.CommandFromStrings(cmd...)
	bytes, _ := cmdArr.Serialize()
	if len(bytes) == 0 {
		slog.Default().Debug("attempt to send empty data to master on handshake", "args", cmd)
	}
	_, err := conn.Write(bytes)
	if err != nil {
		return "", err
	}
	line, err := br.ReadString('\n')
	if err != nil { // exactly one reply line
		return "", fmt.Errorf("handshake %s: %w", cmd[0], err)
	}
	if line[0] == '-' {
		return "", fmt.Errorf("handshake %s rejected: %s", cmd[0], strings.TrimSpace(line))
	}
	slog.Default().Debug("repl -> master", "sent", cmd, "received", line)
	return line, nil
}

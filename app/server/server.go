package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/netip"
	"os"
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
	if err := s.listenMaster(ctx); err != nil {
		logger.Error("Error connecting to master", "err", err.Error())
	}
	logger.Debug("[server] start listening", "addr", s.bind.String(), "port", s.port, "replica", s.env.Repl.Role() != replication.ReplicaRole, "id", s.env.Repl.ReplID())
	defer context.AfterFunc(ctx, func() { l.Close() })()

	for {
		conn, err := l.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				logger.Info("Connection closed")
				os.Exit(0)
			}
			logger.Error("[server][serve] Error accepting connection", "error", err.Error())
			os.Exit(1)
		}

		h := handler.New(s.env, false)
		go h.HandleConnection(ctx, conn)
	}
}

func (s *server) listenMaster(ctx context.Context) error {
	if s.env.Repl.Role() == replication.MasterRole {
		return nil
	}
	slog.Debug("[server][listenMaster] start connecting")
	ctx, cancel := context.WithTimeoutCause(ctx, 5*time.Second, fmt.Errorf("master connection timeout"))
	defer cancel()

	masterAddr, err := s.env.MasterAddr()
	if err != nil {
		return fmt.Errorf("cannot parse master address: %w", err)
	}
	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", masterAddr)
	if err != nil {
		return fmt.Errorf("dial master: %w", err)
	}

	err = s.initMasterConn(conn)
	if err != nil {
		return err
	}

	slog.Debug("[listenMaster] start replication from master")

	h := handler.New(s.env, true)
	go h.HandleConnection(ctx, conn)

	return nil
}

func (s *server) initMasterConn(conn net.Conn) error {
	if _, err := sendToMaster(conn, "PING"); err != nil {
		return err
	}
	if _, err := sendToMaster(conn, "REPLCONF", "listening-port", fmt.Sprint(s.port)); err != nil {
		return err
	}
	if _, err := sendToMaster(conn, "REPLCONF", "capa", "psync2"); err != nil {
		return err
	}
	if _, err := sendToMaster(conn, "PSYNC", "?", "-1"); err != nil {
		return err
	}
	return nil
}

func sendToMaster(conn net.Conn, cmd ...string) ([]byte, error) {
	cmdArr := parser.CommandFromStrings(cmd...)
	bytes, _ := cmdArr.Serialize()
	_, err := conn.Write(bytes)
	if err != nil {
		return nil, err
	}
	var buf [512]byte
	n, err := conn.Read(buf[:])
	return buf[:n], nil
}

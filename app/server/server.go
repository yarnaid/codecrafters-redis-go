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
	if err := s.listenMaster(ctx); err != nil {
		logger.Error("Error connecting to master", "err", err.Error())
	}
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
				logger.Info("Connection closed", "addr", conn.RemoteAddr())
				os.Exit(0)
			}
			logger.Error("Error accepting connection", "error", err.Error())
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
	dailCtx, cancel := context.WithTimeoutCause(ctx, 5*time.Second, fmt.Errorf("master connection timeout"))
	defer cancel()

	masterAddr, err := s.env.MasterAddr()
	if err != nil {
		return fmt.Errorf("cannot parse master address: %w", err)
	}
	var d net.Dialer
	conn, err := d.DialContext(dailCtx, "tcp", masterAddr)
	if err != nil {
		return fmt.Errorf("dial master: %w", err)
	}

	err = s.handshake(conn)
	if err != nil {
		if err1 := conn.Close(); err1 != nil {
			logger.Error("cannot init master conn", "err", err, "err2", err1)
		}
		return err
	}
	_ = conn.SetDeadline(time.Time{})

	slog.Debug("[listenMaster] start replication from master")

	h := handler.New(s.env, true)
	go h.HandleConnection(ctx, conn)

	return nil
}

func (s *server) handshake(conn net.Conn) error {
	br := bufio.NewReader(conn)
	if err := sendToMaster(conn, br, "PING"); err != nil {
		return err
	}
	if err := sendToMaster(conn, br, "REPLCONF", "listening-port", fmt.Sprint(s.port)); err != nil {
		return err
	}
	if err := sendToMaster(conn, br, "REPLCONF", "capa", "psync2"); err != nil {
		return err
	}
	if err := sendToMaster(conn, br, "PSYNC", "?", "-1"); err != nil {
		return err
	}
	hdr, err := br.ReadString('\n')
	if err != nil || hdr[0] != '$' {
		return fmt.Errorf("handshake rdb header: %w", err)
	}
	n, err := strconv.ParseInt(strings.TrimSpace(hdr[1:]), 10, 64)
	if err != nil || n < 0 {
		return errors.New("handshake: bad rdb length")
	}
	_, err = io.CopyN(io.Discard, br, n)
	return err
}

func sendToMaster(conn net.Conn, br *bufio.Reader, cmd ...string) error {
	cmdArr := parser.CommandFromStrings(cmd...)
	bytes, _ := cmdArr.Serialize()
	_, err := conn.Write(bytes)
	if err != nil {
		return err
	}
	line, err := br.ReadString('\n')
	if err != nil { // exactly one reply line
		return fmt.Errorf("handshake %s: %w", cmd[0], err)
	}
	if line[0] == '-' {
		return fmt.Errorf("handshake %s rejected: %s", cmd[0], strings.TrimSpace(line))
	}
	return nil
}

package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"sync"
	"time"

	"my-redis/app/parser"
	"my-redis/app/storage"
)

type server struct {
	cfg        config
	wg         *sync.WaitGroup
	coord      *storage.Coordinator
	state      *ServerState
	masterConn net.Conn
}

func NewServer() (*server, error) {
	coordinator := storage.NewMemoryCoordinator(nil)
	globalWait := sync.WaitGroup{}
	cfg, err := parseConfig(os.Args[1:])
	state := ServerState{id: newServerId()}
	if cfg.replicaof != "" {
		state.role = slaveRole
	} else {
		state.role = masterRole
	}
	if err != nil {
		return nil, err
	}
	return &server{
		cfg:   cfg,
		wg:    &globalWait,
		coord: coordinator,
		state: &state,
	}, nil
}

func (s *server) Serve() {
	l, err := net.Listen("tcp", s.cfg.listenAddr().String())
	if err != nil {
		logger.Error("Failed to bind: %v\ncfg=%c", err.Error(), s.cfg)
		os.Exit(1)
	}
	if err := s.connectToMaster(); err != nil {
		logger.Error("Error connecting to master", "err", err.Error())
	}
	logger.Debug("[server] start listening", "addr", s.cfg.bind.String(), "port", s.cfg.port, "replica", s.cfg.replicaof)
	logger.Debug("[server]", "state", s.state.info())
	for {
		conn, err := l.Accept()
		if err != nil {
			logger.Error("Error accepting connection", "error", err.Error())
			os.Exit(1)
		}
		handler := Handler{Conn: conn, Coordinator: s.coord, TxWait: s.wg, Cfg: &s.cfg, serverState: s.state}
		go handler.HandleConnection()
	}
}

func (s *server) connectToMaster() error {
	if s.cfg.master == nil {
		if s.state.role == slaveRole {
			return fmt.Errorf("cannot start slave without master address")
		}
		return nil
	}
	slog.Debug("[connectToMaster] start connecting")
	ctx, cancel := context.WithTimeoutCause(context.Background(), 5*time.Second, fmt.Errorf("master connection timeout"))
	defer cancel()

	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", s.cfg.master.listenAddr().String())
	if err != nil {
		return fmt.Errorf("dial master: %w", err)
	}
	s.masterConn = conn

	if _, err := s.sendToMaster("PING"); err != nil {
		return err
	}
	if _, err := s.sendToMaster("REPLCONF", "listening-port", fmt.Sprint(s.cfg.port)); err != nil {
		return err
	}
	if _, err := s.sendToMaster("REPLCONF", "capa", "psync2"); err != nil {
		return err
	}
	if _, err := s.sendToMaster("PSYNC", "?", "-1"); err != nil {
		return err
	}

	return err
}

func (s *server) sendToMaster(cmd ...string) ([]byte, error) {
	cmdArr := parser.CommandFromStrings(cmd...)
	bytes, _ := cmdArr.Serialize()
	_, err := s.masterConn.Write(bytes)
	if err != nil {
		return nil, err
	}
	var buf [512]byte
	n, err := s.masterConn.Read(buf[:])
	return buf[:n], nil
}

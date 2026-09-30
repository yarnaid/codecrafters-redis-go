package main

import (
	"net"
	"os"
	"sync"

	"my-redis/app/storage"
)

type server struct {
	cfg   config
	wg    *sync.WaitGroup
	coord *storage.Coordinator
	state *serverState
}

func newServer() (*server, error) {
	coordinator := storage.NewMemoryCoordinator(nil)
	globalWait := sync.WaitGroup{}
	cfg, err := parseConfig(os.Args[1:])
	state := serverState{id: newServerId()}
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

func (s *server) serve() {
	l, err := net.Listen("tcp", s.cfg.listenAddr().String())
	if err != nil {
		logger.Error("Failed to bind: %v\ncfg=%c", err.Error(), s.cfg)
		os.Exit(1)
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

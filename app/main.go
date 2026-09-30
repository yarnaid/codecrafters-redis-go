package main

import (
	"net"
	"os"
	"sync"

	"my-redis/app/storage"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var (
	_ = net.Listen
	_ = os.Exit
)

func main() {
	logger.Info("Logs from your program will appear here!")
	coordinator := storage.NewMemoryCoordinator(nil)
	globalWait := sync.WaitGroup{}
	cfg, err := parseConfig(os.Args[1:])
	if err != nil {
		logger.Error("error parsing config", err)
		return
	}

	l, err := net.Listen("tcp", cfg.listenAddr().String())
	logger.Debug("[server] start listening", "addr", cfg.bind.String(), "port", cfg.port)
	if err != nil {
		logger.Error("Failed to bind: %v\ncfg=%c", err.Error(), cfg)
		os.Exit(1)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			logger.Error("Error accepting connection", "error", err.Error())
			os.Exit(1)
		}
		handler := Handler{Conn: conn, Coordinator: coordinator, TxWait: &globalWait, Cfg: cfg}
		go handler.HandleConnection()
	}
}

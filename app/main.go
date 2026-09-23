package main

import (
	"net"
	"os"

	"my-redis/app/storage"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var (
	_ = net.Listen
	_ = os.Exit
)

var cfg = Config{
	address: "0.0.0.0",
	port:    "6379",
}

func main() {
	logger.Info("Logs from your program will appear here!")
	coordinator := storage.NewMemoryCoordinator(nil)

	l, err := net.Listen("tcp", cfg.address+":"+cfg.port)
	if err != nil {
		logger.Error("Failed to bind", "address", cfg.address, "port", cfg.port, "error", err.Error())
		os.Exit(1)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			logger.Error("Error accepting connection", "error", err.Error())
			os.Exit(1)
		}
		handler := Handler{Conn: conn, Coordinator: coordinator}
		go handler.HandleConnection()
	}
}

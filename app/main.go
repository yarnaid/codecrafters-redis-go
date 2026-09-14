package main

import (
	"net"
	"os"

	"github.com/codecrafters-io/redis-starter-go/app/storage"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

var cfg = Config{
	address: "0.0.0.0",
	port:    "6379",
}

func main() {
	logger.Info("Logs from your program will appear here!")
	var strg = storage.NewStorage()

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
		handler := Handler{conn: conn, strg: strg}
		go handler.handle_connection()
	}
}

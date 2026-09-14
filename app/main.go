package main

import (
	"net"
	"os"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func handle_connection(conn net.Conn) {
	defer func() {
		conn.Close()
		logger.Debug("Connection closed")
	}()

	_, err := conn.Write([]byte(fmt_response("PONG")))
	if err != nil {
		logger.Error("Error writing response", "error", err.Error())
	}

	logger.Info("Handling new connection...")
}

var cfg = Config{
	address: "0.0.0.0",
	port:    "6379",
}

func main() {
	logger.Info("Logs from your program will appear here!")

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
		go handle_connection(conn)
	}
}

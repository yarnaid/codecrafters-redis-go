package main

import (
	"net"
	"os"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var (
	_ = net.Listen
	_ = os.Exit
)

func main() {
	logger.Info("Starting redis")
	server, err := newServer()
	if err != nil {
		logger.Error("cannot start server", "err", err.Error())
		return
	}
	server.serve()
}

package main

import (
	"log/slog"
	"net"
	"os"

	"my-redis/app/server"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var (
	_ = net.Listen
	_ = os.Exit
)

func main() {
	slog.Info("Starting redis")
	server, err := server.NewServer()
	if err != nil {
		slog.Error("cannot start server", "err", err.Error())
		return
	}
	server.Serve()
}

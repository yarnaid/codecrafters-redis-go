package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"my-redis/app/config"
	"my-redis/app/server"
)

func main() {
	slog.Info("Starting redis")
	cfg, err := config.ParseConfig(os.Args[1:])
	if err != nil {
		slog.Error("cannot parse config", "err", err.Error())
	}
	server, err := server.NewServer(cfg)
	if err != nil {
		slog.Error("cannot start server", "err", err.Error())
		return
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	server.Serve(ctx)
}

package handler

import (
	"log/slog"
	"os"

	"my-redis/app/commands"
)

func New(env *commands.Env, fromMaster bool) *handler {
	return &handler{
		logger:     slog.New(slog.NewTextHandler(os.Stderr, nil)),
		env:        env,
		fromMaster: fromMaster,
	}
}

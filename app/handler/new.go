package handler

import (
	"log/slog"

	"my-redis/app/commands"
)

func New(env *commands.Env, fromMaster bool) *handler {
	logger := slog.Default().With("component", "handler")
	if fromMaster {
		logger = logger.With("master", true)
	}
	return &handler{
		logger:     logger,
		env:        env,
		fromMaster: fromMaster,
	}
}

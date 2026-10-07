package handler

import (
	"log/slog"

	"my-redis/app/commands"
)

func New(env *commands.Env, fromMaster bool) *handler {
	logger := slog.Default().With("component", "handler").With("role", env.Repl.Role())
	if fromMaster {
		logger = logger.With("fromMaster", true)
	}
	return &handler{
		logger:     logger,
		env:        env,
		fromMaster: fromMaster,
	}
}

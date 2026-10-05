package replication

import (
	"log/slog"
	"os"
)

func New(role Role, debug bool) *Manager {
	var id replID
	if role == MasterRole {
		id = newReplID()
	} else {
		id = "?"
	}
	var level slog.Level
	if debug {
		level = slog.LevelDebug.Level()
	} else {
		level = slog.LevelInfo.Level()
	}
	return &Manager{
		role:     role,
		replicas: make(map[*Replica]struct{}, 0),
		replID:   id,
		logger:   slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level})).With("component", "replManager"),
	}
}

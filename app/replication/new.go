package replication

import (
	"log/slog"
)

func New(role Role) *Manager {
	var id replID
	if role == MasterRole {
		id = newReplID()
	} else {
		id = "?"
	}
	return &Manager{
		role:     role,
		replicas: make(map[*Replica]struct{}, 0),
		replID:   id,
		logger:   slog.Default().With("component", "replManager"),
	}
}

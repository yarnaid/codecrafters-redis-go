package server

import (
	"sync"

	"my-redis/app/commands"
	"my-redis/app/config"
	"my-redis/app/replication"
	"my-redis/app/storage"
)

func NewServer(cfg *config.Config) (*server, error) {
	var role replication.Role
	if len(cfg.Replicaof) > 0 {
		role = replication.ReplicaRole
	} else {
		role = replication.MasterRole
	}
	coordinator := storage.NewMemoryCoordinator(nil)
	globalWait := sync.WaitGroup{}
	repl := replication.New(role, cfg.Debug)
	env := commands.NewEnv(coordinator, repl, cfg, &globalWait)
	return &server{
		bind: cfg.Bind,
		port: cfg.Port,
		env:  env,
	}, nil
}

package server

import (
	"log/slog"
	"os"
	"path"
	"sync"

	"my-redis/app/commands"
	"my-redis/app/config"
	"my-redis/app/replication"
	"my-redis/app/storage"
	rdbfile "my-redis/app/storage/rdb_file"
)

func NewServer(cfg *config.Config) (*server, error) {
	var role replication.Role
	if len(cfg.Replicaof) > 0 {
		role = replication.ReplicaRole
	} else {
		role = replication.MasterRole
	}
	snap, err := rdbfile.Load(cfg.Dir, cfg.DBFilename)
	if err != nil {
		return nil, err
	}
	slog.Default().Debug("snapshot loaded", "len_keys", len(snap.Data), "data", snap.Data)
	coordinator := storage.NewMemoryCoordinator(snap.ToMemoryBackend())
	globalWait := sync.WaitGroup{}
	repl := replication.New(role)
	env := commands.NewEnv(coordinator, repl, cfg, &globalWait)

	if cfg.AppendOnly == "yes" {
		os.MkdirAll(path.Join(cfg.Dir, cfg.AppendDirName), 0o666)
	}
	return &server{
		bind: cfg.Bind,
		port: cfg.Port,
		env:  env,
	}, nil
}

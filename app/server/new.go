package server

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path"
	"path/filepath"
	"sync"

	"my-redis/app/commands"
	"my-redis/app/config"
	"my-redis/app/handler"
	"my-redis/app/parser"
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
	coordinator := storage.NewMemoryCoordinator(snap.ToMemoryBackend())
	globalWait := sync.WaitGroup{}
	repl := replication.New(role)
	env := commands.NewEnv(coordinator, repl, cfg, &globalWait)

	if cfg.AppendOnly == "yes" {
		aofDir := path.Join(cfg.Dir, cfg.AppendDirName)
		newFileName, err := createAOF(aofDir, cfg.AppendFileName)
		if err != nil {
			if errors.Is(err, os.ErrExist) {
				env.Log.Info("AOF file exists", "name", newFileName)
			} else {
				return nil, err
			}
		}
		createManifest(aofDir, cfg.AppendFileName, newFileName)
		if err != nil {
			if errors.Is(err, os.ErrExist) {
				env.Log.Info("manifest file exists", "name", newFileName)
			} else {
				return nil, err
			}
		}
		manifestName := cfg.AppendFileName + ".manifest"
		m, err := readManifestFile(aofDir, manifestName)
		if err != nil {
			return nil, err
		}
		aofName := filepath.Join(aofDir, m.Name)
		applyAOF(aofName, env)
		env.ManifestChan = make(chan string)
		go startAOFWrite(aofName, env.ManifestChan)
	}

	env.Started = true
	return &server{
		bind: cfg.Bind,
		port: cfg.Port,
		env:  env,
	}, nil
}

func applyAOF(filename string, env *commands.Env) error {
	h := handler.New(env, false)
	sess := commands.NewSession(env.Repl.Role(), nil)
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	commands, err := parser.NewReader(bytes.NewReader(data)).ReadArrays()
	if err != nil {
		return err
	}
	for _, cmd := range commands {
		h.Dispatch(context.Background(), cmd, sess)
	}
	return nil
}

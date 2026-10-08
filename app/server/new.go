package server

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
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
	coordinator := storage.NewMemoryCoordinator(snap.ToMemoryBackend())
	globalWait := sync.WaitGroup{}
	repl := replication.New(role)
	env := commands.NewEnv(coordinator, repl, cfg, &globalWait)

	if cfg.AppendOnly == "yes" {
		aofDir := path.Join(cfg.Dir, cfg.AppendDirName)
		_, err := createAOF(aofDir, cfg.AppendFileName)
		if err != nil {
			return nil, err
		}
		createManifest(aofDir, cfg.AppendFileName)
		if err != nil {
			return nil, err
		}
	}

	return &server{
		bind: cfg.Bind,
		port: cfg.Port,
		env:  env,
	}, nil
}

func createManifest(dir, name string) error {
	f, err := os.OpenFile(filepath.Join(dir, name)+".manifest", os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return err
	}
	if _, err = f.WriteString(fmt.Sprintf("file %s seq 1 type i", name)); err != nil {
		return err
	}
	return nil
}

func createAOF(dir, name string) (string, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}

	newName, err := newAOFFileName(dir, name)
	if err != nil {
		return "", err
	}
	f, err := os.OpenFile(filepath.Join(dir, newName), os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
	if err != nil {
		return "", fmt.Errorf("cannot create file: %w", err)
	}
	if err := f.Close(); err != nil {
		return "", err
	}
	return newName, nil
}

func newAOFFileName(aofPath, name string) (string, error) {
	fileNames, err := filepath.Glob(fmt.Sprintf("%s.*.incr.aof", name))
	if err != nil {
		return "", err
	}
	next, err := getNextAOFIncr(fileNames, name)
	if err != nil {
		return "", nil
	}
	newName := fmt.Sprintf("%s.%d.incr.aof", name, next)
	return newName, nil
}

func getNextAOFIncr(prev []string, fname string) (int, error) {
	re, err := regexp.Compile(fmt.Sprintf("^%s\\.(\\d+)\\.incr.aof$", fname))
	if err != nil {
		return -1, err
	}
	res := 0
	for _, s := range prev {
		sm := re.FindStringSubmatch(s)
		if len(sm) == 0 {
			continue
		}
		num, err := strconv.Atoi(sm[1])
		if err != nil {
			return -1, err
		}
		res = max(res, num)
	}
	return res + 1, nil
}

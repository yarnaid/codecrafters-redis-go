package commands

import (
	"fmt"
	"log/slog"
	"net/netip"
	"strings"
	"sync"

	"my-redis/app/config"
	"my-redis/app/replication"
	"my-redis/app/storage"
)

type Env struct {
	Repl         *replication.Manager
	Coordinator  *storage.Coordinator
	Reg          *Registry
	Log          *slog.Logger
	Cfg          *config.Config
	Wg           *sync.WaitGroup
	ManifestChan chan string
}

func NewEnv(coord *storage.Coordinator, repl *replication.Manager, cfg *config.Config, wg *sync.WaitGroup) *Env {
	return &Env{
		Repl:        repl,
		Coordinator: coord,
		Reg:         DefaultRegistry(),
		Log:         slog.Default(),
		Cfg:         cfg,
		Wg:          wg,
	}
}

func (e *Env) MasterAddr() (string, error) {
	s := strings.Split(e.Cfg.Replicaof, " ")
	pp := config.Port(0)
	p := &pp
	err := p.UnmarshalText([]byte(s[1]))
	if err != nil {
		return "", fmt.Errorf("parsing master addr bytes: %w", err)
	}
	if s[0] == "localhost" {
		s[0] = "127.0.0.1"
	}
	a, err := netip.ParseAddr(s[0])
	if err != nil {
		return "", fmt.Errorf("parsing master addr to ip: %w", err)
	}
	return netip.AddrPortFrom(a, uint16(*p)).String(), nil
}

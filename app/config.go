package main

import (
	"errors"
	"flag"
	"fmt"
	"net/netip"
	"strconv"
	"strings"

	"github.com/peterbourgon/ff/fftoml"
	"github.com/peterbourgon/ff/v3"
)

type port uint16

func (p port) MarshalText() ([]byte, error) {
	return strconv.AppendUint(nil, uint64(p), 10), nil
}

func (p *port) UnmarshalText(b []byte) error {
	v, err := strconv.ParseUint(string(b), 10, 16) // bitSize 16 enforces 0-65535
	if err != nil {
		return fmt.Errorf("invalid port %q: want 0-65535", b)
	}
	*p = port(v)
	return nil
}

type master struct {
	bind netip.Addr
	port port
}

func (c *master) listenAddr() netip.AddrPort {
	return netip.AddrPortFrom(c.bind, uint16(c.port))
}

type config struct {
	bind      netip.Addr
	port      port
	replicaof string
	master    *master
}

func (c config) listenAddr() netip.AddrPort {
	return netip.AddrPortFrom(c.bind, uint16(c.port))
}

func (c config) validate() error {
	if !c.bind.IsValid() {
		return errors.New("bind: empty or invalid address")
	}
	return nil
}

func parseConfig(args []string) (config, error) {
	var cfg config

	fs := flag.NewFlagSet("redis", flag.ContinueOnError)
	fs.String("config", "", "TOML config file (optional)")
	fs.StringVar(&cfg.replicaof, "replicaof", "", "master server for replication")
	fs.TextVar(&cfg.bind, "bind", netip.MustParseAddr("127.0.0.1"), "bind IP address")
	fs.TextVar(&cfg.port, "port", port(6379), "listen port")
	if err := ff.Parse(
		fs, args,
		ff.WithConfigFileFlag("config"),
		ff.WithConfigFileParser(fftoml.Parser),
		ff.WithEnvVarPrefix("REDIS"),
	); err != nil {
		return cfg, err
	}

	if cfg.replicaof != "" {
		s := strings.Split(cfg.replicaof, " ")
		pp := port(0)
		p := &pp
		err := p.UnmarshalText([]byte(s[1]))
		if err != nil {
			return cfg, fmt.Errorf("parsing master addr bytes: %w", err)
		}
		if s[0] == "localhost" {
			s[0] = "127.0.0.1"
		}
		a, err := netip.ParseAddr(s[0])
		if err != nil {
			return cfg, fmt.Errorf("parsing master addr to ip: %w", err)
		}
		master := master{
			bind: a,
			port: *p,
		}
		cfg.master = &master
	}
	cfg.bind = cfg.bind.Unmap()
	err := cfg.validate()
	return cfg, err
}

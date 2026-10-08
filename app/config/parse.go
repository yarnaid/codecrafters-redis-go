package config

import (
	"errors"
	"flag"
	"net/netip"

	"github.com/peterbourgon/ff/fftoml"
	"github.com/peterbourgon/ff/v3"
)

func (c Config) validate() error {
	if !c.Bind.IsValid() {
		return errors.New("bind: empty or invalid address")
	}
	return nil
}

func ParseConfig(args []string) (*Config, error) {
	var cfg Config

	fs := flag.NewFlagSet("redis", flag.ContinueOnError)
	fs.String("config", "", "TOML config file (optional)")
	fs.StringVar(&cfg.Replicaof, "replicaof", "", "master server for replication")
	fs.StringVar(&cfg.Dir, "dir", "app/storage/rdb_file/testdata/", "db file location")
	fs.StringVar(&cfg.DBFilename, "dbfilename", "dump_test.rdb", "db file name")
	fs.TextVar(&cfg.Bind, "bind", netip.MustParseAddr("127.0.0.1"), "bind IP address")
	fs.TextVar(&cfg.Port, "port", Port(6379), "listen port")
	fs.BoolVar(&cfg.Debug, "debug", true, "enable debug mode")
	if err := ff.Parse(
		fs, args,
		ff.WithConfigFileFlag("config"),
		ff.WithConfigFileParser(fftoml.Parser),
		ff.WithEnvVarPrefix("REDIS"),
	); err != nil {
		return nil, err
	}

	cfg.Bind = cfg.Bind.Unmap()
	err := cfg.validate()
	return &cfg, err
}

package config

import (
	"errors"
	"flag"
	"net/netip"
	"os"

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
	cwd, _ := os.Getwd()

	fs := flag.NewFlagSet("redis", flag.ContinueOnError)
	fs.String("config", "", "TOML config file (optional)")
	fs.StringVar(&cfg.Replicaof, "replicaof", "", "master server for replication")

	fs.StringVar(&cfg.Dir, "dir", cwd, "db file location")
	fs.StringVar(&cfg.DBFilename, "dbfilename", "dump_test.rdb", "db file name")

	fs.StringVar(&cfg.AppendOnly, "appendonly", "no", "append only yes/no")
	fs.StringVar(&cfg.AppendDirName, "appenddirname", "appendonlydir", "append only dir name")
	fs.StringVar(&cfg.AppendFileName, "appendfilename", "appendonly.aof", "append only file name")
	fs.StringVar(&cfg.AppendFSync, "appendfsync", "everysec", "append only sync freq")

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

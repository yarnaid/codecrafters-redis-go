package commands

import (
	"context"
	"fmt"
	"strings"

	"my-redis/app/parser"
)

type ConfigCommand struct {
	Subcmd string
	Args   []string
}

var _ Command = (*ConfigCommand)(nil)

func (c *ConfigCommand) Execute(ctx context.Context, env *Env, sess *Session) (parser.Serializable, error) {
	var v string
	cfg, _ := env.Cfg.ToMap()
	v, ok := cfg[c.Args[0]]
	if !ok {
		return nil, fmt.Errorf("config key not found: %s", c.Args[0])
	}
	return parser.CommandFromStrings(c.Args[0], v), nil
}

func parseConfig(args []string) (Command, error) {
	switch strings.ToLower(args[0]) {
	case "get":
		return &ConfigCommand{"get", args[1:]}, nil
	}
	return nil, fmt.Errorf("cannot parse config command args: %v", args)
}

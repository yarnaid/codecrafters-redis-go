package commands

import (
	"context"
	"fmt"

	"my-redis/app/parser"
)

type GeoDistCommand struct {
	Name   string
	L1, L2 string
}

var _ Command = (*GeoDistCommand)(nil)

func (g *GeoDistCommand) Execute(ctx context.Context, env *Env, sess *Session) (parser.Serializable, error) {
	res, err := env.Coordinator.GeoDist(g.Name, g.L1, g.L2)
	if err != nil {
		return nil, err
	}
	return parser.BulkString(fmt.Sprint(res)), nil
}

func parseGeoDist(args []string) (Command, error) {
	return &GeoDistCommand{args[0], args[1], args[2]}, nil
}

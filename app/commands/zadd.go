package commands

import (
	"context"
	"fmt"
	"strconv"

	"my-redis/app/parser"
)

type ZAddCommand struct {
	Key   string
	Score float64
	Name  string
}

var _ Command = (*ZAddCommand)(nil)

func (z ZAddCommand) Execute(ctx context.Context, env *Env, sess *Session) (parser.Serializable, error) {
	n, err := env.Coordinator.ZAdd(z.Key, z.Name, z.Score)
	return parser.Int(n), err
}

func parseZAdd(args []string) (Command, error) {
	score, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		return nil, fmt.Errorf("cannot parse zadd scores: %w", err)
	}
	return ZAddCommand{args[0], score, args[2]}, nil
}

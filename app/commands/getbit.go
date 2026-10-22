package commands

import (
	"context"
	"strconv"

	"my-redis/app/parser"
)

type GetBitCommand struct {
	Key    string
	Offset int
}

var _ Command = (*GetBitCommand)(nil)

func (g *GetBitCommand) Execute(ctx context.Context, env *Env, sess *Session) (parser.Serializable, error) {
	res, err := env.Coordinator.GetBit(g.Key, g.Offset)
	if err != nil {
		return nil, err
	}
	return parser.Int(res), nil
}

func parseGetBit(args []string) (Command, error) {
	offset, err := strconv.Atoi(args[1])
	if err != nil {
		return nil, err
	}
	return &GetBitCommand{args[0], offset}, nil
}

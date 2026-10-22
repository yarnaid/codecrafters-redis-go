package commands

import (
	"context"

	"my-redis/app/parser"
)

type BitOpCommand struct {
	Op     string
	Dest   string
	K1, K2 string
}

var _ Command = (*BitOpCommand)(nil)

func (b *BitOpCommand) Execute(ctx context.Context, env *Env, sess *Session) (parser.Serializable, error) {
	res, err := env.Coordinator.BitOp(b.Op, b.Dest, b.K1, b.K2)
	if err != nil {
		return nil, err
	}
	return parser.Int(res), nil
}

func parseBitOp(args []string) (Command, error) {
	return &BitOpCommand{args[0], args[1], args[2], args[3]}, nil
}

package commands

import (
	"context"

	"my-redis/app/parser"
)

type LLen struct {
	Key string
}

var _ Command = (*LLen)(nil)

func (l LLen) Execute(ctx context.Context, env *Env, _ *Session) (parser.Serializable, error) {
	return parser.Int(env.Coordinator.LLen(l.Key)), nil
}

func parseLLen(args []string) (Command, error) {
	return LLen{args[0]}, nil
}

package commands

import (
	"context"

	"my-redis/app/parser"
)

type LPushCommand struct {
	Key    string
	Values []parser.Serializable
}

var _ Command = (*LPushCommand)(nil)

func (r LPushCommand) Execute(ctx context.Context, env *Env, _ *Session) (parser.Serializable, error) {
	return parser.Int(env.Coordinator.Prepend(r.Key, r.Values...)), nil
}

func parseLPush(args []string) (Command, error) {
	return LPushCommand{args[0], parser.ToSerializableSlice(args[1:])}, nil
}

package commands

import (
	"context"
	"log/slog"

	"my-redis/app/parser"
)

type RPushCommand struct {
	Key    string
	Values []parser.Serializable
}

var _ Command = (*RPushCommand)(nil)

func (r RPushCommand) Execute(ctx context.Context, env *Env, _ *Session) (parser.Serializable, error) {
	slog.Debug("[Command][RPUSH]", "key", r.Key, "values", r.Values)
	return parser.Int(env.Coordinator.Append(r.Key, r.Values...)), nil
}

func parseRPush(args []string) (Command, error) {
	return RPushCommand{args[0], parser.ToSerializableSlice(args[1:])}, nil
}

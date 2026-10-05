package commands

import (
	"context"

	"my-redis/app/parser"
)

type TypeCommand struct {
	Key string
}

var _ Command = (*TypeCommand)(nil)

func (t TypeCommand) Execute(ctx context.Context, env *Env, _ *Session) (parser.Serializable, error) {
	k, err := env.Coordinator.Type(t.Key)
	if err != nil {
		return nil, err
	}
	s := string(k)
	return parser.SimpleString(s), nil
}

func parseType(args []string) (Command, error) {
	return TypeCommand{Key: args[0]}, nil
}

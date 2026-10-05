package commands

import (
	"context"

	"my-redis/app/parser"
)

type MultiCommand struct{}

var _ Command = (*MultiCommand)(nil)

func (e MultiCommand) Execute(ctx context.Context, env *Env, sess *Session) (parser.Serializable, error) {
	return parser.SimpleString("OK"), nil
}

func parseMulti(args []string) (Command, error) {
	return MultiCommand{}, nil
}

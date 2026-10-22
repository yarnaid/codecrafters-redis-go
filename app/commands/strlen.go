package commands

import (
	"context"

	"my-redis/app/parser"
)

type StrLenCommand struct {
	Key string
}

var _ Command = (*StrLenCommand)(nil)

func (s *StrLenCommand) Execute(ctx context.Context, env *Env, sess *Session) (parser.Serializable, error) {
	res, _ := env.Coordinator.StrLen(s.Key)
	return parser.Int(res), nil
}

func parseStrLen(args []string) (Command, error) {
	return &StrLenCommand{args[0]}, nil
}

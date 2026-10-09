package commands

import (
	"context"

	"my-redis/app/parser"
)

type SubscribeCommand struct {
	Subs []string
}

var _ Command = (*SubscribeCommand)(nil)

func (s *SubscribeCommand) Execute(ctx context.Context, env *Env, sess *Session) (parser.Serializable, error) {
	return parser.Array[parser.Serializable]{parser.BulkString("subscribe"), parser.BulkString(s.Subs[0]), parser.Int(1)}, nil
}

func parseSubscribe(args []string) (Command, error) {
	return &SubscribeCommand{args}, nil
}

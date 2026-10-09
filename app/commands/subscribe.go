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
	n := env.Broker.Subscribe(sess.Client, s.Subs[0])
	sess.Mode = ModeSub
	return parser.Array[parser.Serializable]{parser.BulkString("subscribe"), parser.BulkString(s.Subs[0]), parser.Int(n)}, nil
}

func parseSubscribe(args []string) (Command, error) {
	return &SubscribeCommand{args}, nil
}

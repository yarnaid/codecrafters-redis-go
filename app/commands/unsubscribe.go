package commands

import (
	"context"

	"my-redis/app/parser"
)

type UnsubscribeCommand struct {
	Subs []string
}

var _ Command = (*UnsubscribeCommand)(nil)

func (s *UnsubscribeCommand) Execute(ctx context.Context, env *Env, sess *Session) (parser.Serializable, error) {
	n := env.Broker.Unsubscribe(sess.Client, s.Subs[0])
	sess.Mode = ModeSub
	return parser.Array[parser.Serializable]{parser.BulkString("unsubscribe"), parser.BulkString(s.Subs[0]), parser.Int(n)}, nil
}

func parseUnsubscribe(args []string) (Command, error) {
	return &UnsubscribeCommand{args}, nil
}

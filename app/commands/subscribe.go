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
	if _, ok := sess.Subscriptions[s.Subs[0]]; !ok {
		sess.Subscriptions[s.Subs[0]] = NewSubscription()
	}
	sess.Mode = ModeSub
	return parser.Array[parser.Serializable]{parser.BulkString("subscribe"), parser.BulkString(s.Subs[0]), parser.Int(len(sess.Subscriptions))}, nil
}

func parseSubscribe(args []string) (Command, error) {
	return &SubscribeCommand{args}, nil
}

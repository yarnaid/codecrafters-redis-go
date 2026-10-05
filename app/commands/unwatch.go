package commands

import (
	"context"

	"my-redis/app/parser"
)

type UnwatchCommand struct{}

var _ Command = (*UnwatchCommand)(nil)

func (w UnwatchCommand) Execute(ctx context.Context, env *Env, sess *Session) (parser.Serializable, error) {
	sess.WatchList = sess.WatchList[:0]
	return parser.SimpleString("OK"), nil
}

func parseUnwatch(args []string) (Command, error) {
	return UnwatchCommand{}, nil
}

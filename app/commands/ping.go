package commands

import (
	"context"

	"my-redis/app/parser"
)

type PingCommand struct{}

var _ Command = (*PingCommand)(nil)

func (p PingCommand) Execute(ctx context.Context, _ *Env, sess *Session) (parser.Serializable, error) {
	if sess.Mode == ModeSub {
		return parser.Array[parser.BulkString]{parser.BulkString("pong"), parser.BulkString("")}, nil
	}
	return parser.SimpleString("PONG"), nil
}

func parsePing(args []string) (Command, error) {
	return PingCommand{}, nil
}

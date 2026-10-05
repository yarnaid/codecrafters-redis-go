package commands

import (
	"context"

	"my-redis/app/parser"
)

type PingCommand struct{}

var _ Command = (*PingCommand)(nil)

func (p PingCommand) Execute(ctx context.Context, _ *Env, _ *Session) (parser.Serializable, error) {
	return parser.SimpleString("PONG"), nil
}

func parsePing(args []string) (Command, error) {
	return PingCommand{}, nil
}

package commands

import (
	"context"

	"my-redis/app/parser"
)

type EchoCommand struct {
	Msg string
}

var _ Command = (*EchoCommand)(nil)

func (e EchoCommand) Execute(_ context.Context, env *Env, _ *Session) (parser.Serializable, error) {
	return parser.BulkString(e.Msg), nil
}

func parseEcho(args []string) (Command, error) {
	return EchoCommand{Msg: args[0]}, nil
}

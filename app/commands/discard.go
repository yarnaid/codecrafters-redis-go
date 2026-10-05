package commands

import (
	"context"

	"my-redis/app/parser"
)

type DiscardCommand struct {
	TxStarted bool
}

var _ Command = (*DiscardCommand)(nil)

func (e DiscardCommand) Execute(_ context.Context, _ *Env, sess *Session) (parser.Serializable, error) {
	return parser.SimpleString("OK"), nil
}

func parseDiscard(args []string) (Command, error) {
	return DiscardCommand{}, nil
}

package commands

import (
	"context"

	"my-redis/app/parser"
)

type ACLCommand struct{}

var _ Command = (*ACLCommand)(nil)

func (a *ACLCommand) Execute(ctx context.Context, env *Env, sess *Session) (parser.Serializable, error) {
	return parser.BulkString("default"), nil
}

func parseACL(args []string) (Command, error) {
	return &ACLCommand{}, nil
}

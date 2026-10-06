// Package commands ...
package commands

import (
	"context"

	"my-redis/app/parser"
)

type Command interface {
	Execute(ctx context.Context, env *Env, sess *Session) (parser.Serializable, error)
}

type CommandCommand struct{}

var _ Command = (*CommandCommand)(nil)

func (c CommandCommand) Execute(ctx context.Context, env *Env, sess *Session) (parser.Serializable, error) {
	return parser.SimpleString("OK"), nil
}

func parseCommand(args []string) (Command, error) {
	return CommandCommand{}, nil
}

package commands

import (
	"context"

	"my-redis/app/parser"
)

type PSyncCommand struct{}

var _ Command = (*PSyncCommand)(nil)

func (p PSyncCommand) Execute(ctx context.Context, env *Env, sess *Session) (parser.Serializable, error) {
	if sess.Register != nil {
		sess.Register()
	}
	return nil, nil
}

func parsePSync(args []string) (Command, error) {
	return PSyncCommand{}, nil
}

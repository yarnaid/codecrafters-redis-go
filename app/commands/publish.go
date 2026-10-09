package commands

import (
	"context"

	"my-redis/app/parser"
)

type PublishCommand struct {
	Chan string
	Msg  string
}

var _ Command = (*PublishCommand)(nil)

func (p PublishCommand) Execute(ctx context.Context, env *Env, sess *Session) (parser.Serializable, error) {
	n := env.Broker.Publish(p.Chan, p.Msg)
	return parser.Int(n), nil
}

func parsePublish(args []string) (Command, error) {
	return PublishCommand{Chan: args[0], Msg: args[1]}, nil
}

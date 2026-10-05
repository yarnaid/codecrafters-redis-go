package commands

import (
	"context"

	"my-redis/app/parser"
	"my-redis/app/storage"
)

type IncrCommand struct {
	Key string
}

var _ Command = (*IncrCommand)(nil)

func (i IncrCommand) Execute(ctx context.Context, env *Env, sess *Session) (parser.Serializable, error) {
	v, err := env.Coordinator.Incr(i.Key)
	if err != nil {
		switch err.(type) {
		case *storage.WrongTypeError:
			return parser.SimpleError("ERR value is not an integer or out of range"), nil
		default:
			return nil, err
		}
	}
	return v, nil
}

func parseIncr(args []string) (Command, error) {
	return IncrCommand{Key: args[0]}, nil
}

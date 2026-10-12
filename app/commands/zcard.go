package commands

import (
	"context"
	"errors"

	"my-redis/app/parser"
	"my-redis/app/storage"
)

type ZCardCommand struct {
	Key string
}

var _ Command = (*ZCardCommand)(nil)

func (z *ZCardCommand) Execute(ctx context.Context, env *Env, sess *Session) (parser.Serializable, error) {
	res, err := env.Coordinator.ZCard(z.Key)
	if err != nil {
		if errors.Is(err, &storage.NotFoundError{}) {
			return parser.Int(0), nil
		}
		return nil, err
	}
	return parser.Int(res), err
}

func parseZCard(args []string) (Command, error) {
	return &ZCardCommand{args[0]}, nil
}

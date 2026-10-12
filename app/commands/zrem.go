package commands

import (
	"context"
	"errors"

	"my-redis/app/parser"
	"my-redis/app/storage"
)

type ZRem struct {
	Key  string
	Name string
}

var _ Command = (*ZRem)(nil)

func (z ZRem) Execute(ctx context.Context, env *Env, sess *Session) (parser.Serializable, error) {
	res, err := env.Coordinator.ZRem(z.Key, z.Name)
	if err != nil {
		if errors.Is(err, &storage.NotFoundError{}) {
			return parser.Int(0), nil
		}
		return nil, err
	}
	return parser.Int(res), nil
}

func parseZRem(args []string) (Command, error) {
	return ZRem{args[0], args[1]}, nil
}

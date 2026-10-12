package commands

import (
	"context"
	"errors"
	"fmt"

	"my-redis/app/parser"
	"my-redis/app/storage"
)

type ZScoreCommand struct {
	Key  string
	Name string
}

var _ Command = (*ZScoreCommand)(nil)

func (z ZScoreCommand) Execute(ctx context.Context, env *Env, sess *Session) (parser.Serializable, error) {
	rank, err := env.Coordinator.ZScore(z.Key, z.Name)
	if err != nil {
		if errors.Is(err, &storage.NotFoundError{}) {
			return parser.BulkNullString(""), nil
		}
		return nil, err
	}
	return parser.BulkString(fmt.Sprint(rank)), nil
}

func parseZScore(args []string) (Command, error) {
	return ZScoreCommand{args[0], args[1]}, nil
}

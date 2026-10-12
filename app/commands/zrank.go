package commands

import (
	"context"
	"errors"

	"my-redis/app/parser"
	"my-redis/app/storage"
)

type ZRankCommand struct {
	Key  string
	Name string
}

var _ Command = (*ZRankCommand)(nil)

func (z *ZRankCommand) Execute(ctx context.Context, env *Env, sess *Session) (parser.Serializable, error) {
	rank, err := env.Coordinator.ZRank(z.Key, z.Name)
	if err != nil {
		if errors.Is(err, &storage.NotFoundError{}) {
			return parser.BulkNullString(""), nil
		}
		return nil, err
	}
	return parser.Int(rank), nil
}

func parseZRank(args []string) (Command, error) {
	return &ZRankCommand{args[0], args[1]}, nil
}

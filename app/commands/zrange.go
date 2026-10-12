package commands

import (
	"context"
	"errors"
	"strconv"

	"my-redis/app/parser"
	"my-redis/app/storage"
)

type ZRangeCommand struct {
	Key        string
	Start, End int
}

var _ Command = (*ZRangeCommand)(nil)

func (z *ZRangeCommand) Execute(ctx context.Context, env *Env, sess *Session) (parser.Serializable, error) {
	res, err := env.Coordinator.ZRange(z.Key, z.Start, z.End)
	if err != nil {
		if errors.Is(err, &storage.NotFoundError{}) {
			return parser.NullArray(0), nil
		}
		return nil, err
	}
	return parser.StringsToArray(res...), nil
}

func parseZRange(args []string) (Command, error) {
	start, err := strconv.Atoi(args[1])
	if err != nil {
		return nil, err
	}

	end, err := strconv.Atoi(args[2])
	if err != nil {
		return nil, err
	}
	return &ZRangeCommand{args[0], start, end}, nil
}

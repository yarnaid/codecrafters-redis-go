package commands

import (
	"context"

	"my-redis/app/parser"
)

type XRangeCommand struct {
	Key        string
	Start, End string
}

var _ Command = (*XRangeCommand)(nil)

func (l XRangeCommand) Execute(ctx context.Context, env *Env, sess *Session) (parser.Serializable, error) {
	slice, err := env.Coordinator.XRange(l.Key, l.Start, l.End)
	if err != nil {
		return nil, err
	}
	res := make(parser.Array[parser.Serializable], len(slice))
	for i, v := range slice {
		res[i] = v
	}
	return res, nil
}

func parseXRange(args []string) (Command, error) {
	return XRangeCommand{Key: args[0], Start: args[1], End: args[2]}, nil
}

package commands

import (
	"context"
	"strconv"

	"my-redis/app/parser"
)

type BitCountCommand struct {
	Key        string
	Start, End int
}

var _ Command = (*BitCountCommand)(nil)

func (b *BitCountCommand) Execute(ctx context.Context, env *Env, sess *Session) (parser.Serializable, error) {
	res, ok := env.Coordinator.BitCount(b.Key, b.Start, b.End)
	if !ok {
		return parser.Int(0), nil
	}
	return parser.Int(res), nil
}

func parseBitCount(args []string) (Command, error) {
	switch len(args) {
	case 1:
		return &BitCountCommand{args[0], 0, -1}, nil
	case 3:
		start, err := strconv.Atoi(args[1])
		if err != nil {
			return nil, err
		}
		end, err := strconv.Atoi(args[2])
		if err != nil {
			return nil, err
		}
		return &BitCountCommand{args[0], start, end}, nil
	}
	return nil, &InvalidCommandError{"BITCOUNT", args}
}

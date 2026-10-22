package commands

import (
	"context"
	"strconv"

	"my-redis/app/parser"
)

type SetBitCommand struct {
	Key    string
	Offset int
	Zero   bool
}

var _ Command = (*SetBitCommand)(nil)

func (s *SetBitCommand) Execute(ctx context.Context, env *Env, sess *Session) (parser.Serializable, error) {
	res, err := env.Coordinator.SetBit(s.Key, s.Offset, s.Zero)
	if err != nil {
		return nil, err
	}
	return parser.Int(res), nil
}

func parseSetBit(args []string) (Command, error) {
	var zero bool
	switch args[2] {
	case "0":
		zero = true
	case "1":
	default:
		return nil, &InvalidCommandError{"setbit", args}
	}
	offset, err := strconv.Atoi(args[1])
	if err != nil {
		return nil, err
	}
	return &SetBitCommand{args[0], offset, zero}, nil
}

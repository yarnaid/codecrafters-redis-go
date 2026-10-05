package commands

import (
	"context"
	"fmt"
	"strconv"

	"my-redis/app/parser"
)

type LRangeCommand struct {
	Key        string
	Start, End int
}

var _ Command = (*LRangeCommand)(nil)

func (l LRangeCommand) Execute(ctx context.Context, env *Env, _ *Session) (parser.Serializable, error) {
	return env.Coordinator.LRange(l.Key, l.Start, l.End), nil
}

func parseLRange(args []string) (Command, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("LRANGE args format: key, start, end")
	}
	start, err := strconv.Atoi(args[1])
	if err != nil {
		return nil, err
	}
	end, err := strconv.Atoi(args[2])
	if err != nil {
		return nil, err
	}
	return LRangeCommand{args[0], start, end}, nil
}

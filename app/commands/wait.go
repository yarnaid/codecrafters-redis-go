package commands

import (
	"context"
	"strconv"
	"time"

	"my-redis/app/parser"
)

type WaitCommand struct {
	N        int
	Duration time.Duration
}

var _ Command = (*WaitCommand)(nil)

func (w WaitCommand) Execute(ctx context.Context, env *Env, sess *Session) (parser.Serializable, error) {
	count := env.Repl.Count()
	return parser.Int(count), nil
}

func parseWait(args []string) (Command, error) {
	n := -1
	t := time.Duration(0)
	var err error
	if len(args) >= 1 {
		n, err = strconv.Atoi(args[0])
		if err != nil {
			return nil, err
		}
	}
	if len(args) >= 2 {
		tt, err := strconv.Atoi(args[1])
		if err != nil {
			return nil, err
		}
		t = time.Duration(tt)
	}
	return WaitCommand{N: n, Duration: t}, nil
}

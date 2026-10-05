package commands

import (
	"context"
	"fmt"
	"strconv"

	"my-redis/app/parser"
)

type LPop struct {
	Key string
	N   int
}

var _ Command = (*LPop)(nil)

func (l LPop) Execute(ctx context.Context, env *Env, _ *Session) (parser.Serializable, error) {
	v, err := env.Coordinator.LPop(l.Key, l.N)
	var vv parser.Serializable
	if len(v) == 1 {
		vv = parser.ToSerializable(v[0])
	} else {
		vv = parser.ToSerializable(v)
	}
	return vv, err
}

func parseLPop(args []string) (Command, error) {
	if len(args) == 1 {
		return LPop{args[0], 1}, nil
	} else if len(args) == 2 {
		n, err := strconv.Atoi(args[1])
		if err != nil {
			return nil, err
		}
		return LPop{args[0], n}, nil
	}
	return nil, fmt.Errorf("wrong args number %d, must be 1 only", len(args)-1)
}

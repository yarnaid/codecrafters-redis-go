package commands

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"my-redis/app/parser"
	"my-redis/app/storage"
)

type SetCommand struct {
	Key   string
	Value interface{}
	TTL   time.Duration
}

var _ Command = (*SetCommand)(nil)

func (s SetCommand) Execute(ctx context.Context, env *Env, _ *Session) (parser.Serializable, error) {
	ok := env.Coordinator.Set(s.Key, s.Value, s.TTL, storage.KindString)
	if !ok {
		return nil, fmt.Errorf("cannot set `%v` to key `%v`", s.Value, s.Key)
	}
	return parser.SimpleString("OK"), nil
}

func parseSet(args []string) (Command, error) {
	switch len(args) {
	case 2:
		return SetCommand{Key: args[0], Value: args[1]}, nil
	case 4:
		if strings.ToUpper(args[2]) == "PX" {
			ttlMs, err := strconv.Atoi(args[3])
			if err != nil {
				return nil, fmt.Errorf("cannot candle set command TTL: %w", err)
			}
			return SetCommand{Key: args[0], Value: args[1], TTL: time.Duration(ttlMs) * time.Millisecond}, nil
		}
		return nil, fmt.Errorf("SET got wrong arg: %v", args[2])
	default:
		return nil, fmt.Errorf("SET got wrong args number: %d", len(args))
	}
}

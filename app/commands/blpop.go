package commands

import (
	"context"
	"log/slog"
	"math"
	"strconv"
	"time"

	"my-redis/app/parser"
)

type BLPop struct {
	Key     string
	Timeout time.Duration
}

var _ Command = (*BLPop)(nil)

func (l BLPop) Execute(_ context.Context, env *Env, _ *Session) (parser.Serializable, error) {
	slog.Debug("[Command][BLPOP]", "key", l.Key, "timeout", l.Timeout)
	v, err := env.Coordinator.BLPop(l.Key, l.Timeout)
	if err != nil {
		slog.Error("[Command][BLPOP] failed", "key", l.Key, "timeout", l.Timeout, "err", err)
		return parser.NullArray(0), nil
	}
	vv := parser.ToSerializable(v)
	res := parser.Array[parser.Serializable]{parser.BulkString(l.Key), vv}
	slog.Debug("[Command][BLPOP] finished", "key", l.Key, "timeout", l.Timeout, "res", res)
	return res, nil
}

func parseBLPop(args []string) (Command, error) {
	timeout, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		return nil, err
	}
	return BLPop{Key: args[0], Timeout: time.Duration(math.Round(timeout * float64(time.Second)))}, nil
}

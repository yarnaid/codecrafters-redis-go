package commands

import (
	"log/slog"
	"time"

	"my-redis/app/parser"
	"my-redis/app/storage"
)

type BLPop struct {
	S       *storage.Coordinator
	Key     string
	Timeout time.Duration
}

func (l BLPop) Execute() (parser.Serializable, error) {
	slog.Debug("[Command][BLPOP]", "key", l.Key, "timeout", l.Timeout)
	v, err := l.S.BLPop(l.Key, l.Timeout)
	if err != nil {
		slog.Error("[Command][BLPOP] failed", "key", l.Key, "timeout", l.Timeout, "err", err)
		return parser.NullArray(0), nil
	}
	vv := parser.ToSerializable(v)
	res := parser.Array[parser.Serializable]{parser.BulkString(l.Key), vv}
	slog.Debug("[Command][BLPOP] finished", "key", l.Key, "timeout", l.Timeout, "res", res)
	return res, nil
}

func (l BLPop) Validate() error {
	return nil
}

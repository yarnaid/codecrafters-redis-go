package commands

import (
	"context"
	"fmt"
	"strconv"

	"my-redis/app/parser"
)

type GetCommand struct {
	Key string
}

var _ Command = (*GetCommand)(nil)

func (c GetCommand) Execute(ctx context.Context, env *Env, _ *Session) (parser.Serializable, error) {
	val, ok := env.Coordinator.Get(c.Key)
	if !ok {
		return parser.BulkNullString(""), nil
	}
	switch val := val.(type) {
	case string:
		return parser.BulkString(val), nil
	case int:
		// slog.Debug("[Command][GET]", "int", val)
		return parser.BulkString(strconv.FormatInt(int64(val), 10)), nil
	case parser.BulkString:
		return val, nil
	default:
		return nil, fmt.Errorf("unsupported type %T", val)
	}
}

func parseGet(args []string) (Command, error) {
	return GetCommand{Key: args[0]}, nil
}

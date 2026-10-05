package commands

import (
	"context"

	"my-redis/app/parser"
	"my-redis/app/storage"
)

type XaddCommand struct {
	Key    string
	Id     string
	Values []string
}

var _ Command = (*XaddCommand)(nil)

func (s XaddCommand) Execute(ctx context.Context, env *Env, _ *Session) (parser.Serializable, error) {
	newKey, err := env.Coordinator.Xadd(s.Key, s.Id, s.Values...)
	if err != nil {
		switch err.(type) {
		case storage.InvalidStreamIdSeq:
			return parser.SimpleError("ERR The ID specified in XADD is equal or smaller than the target stream top item"), nil
		case storage.ZeroStreamId:
			return parser.SimpleError("ERR The ID specified in XADD must be greater than 0-0"), nil
		default:
			return nil, err
		}
	}
	return parser.BulkString(newKey), nil
}

func parseXAdd(args []string) (Command, error) {
	return XaddCommand{Key: args[0], Id: args[1], Values: args[2:]}, nil
}

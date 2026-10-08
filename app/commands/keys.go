package commands

import (
	"context"

	"my-redis/app/parser"
)

type KeysCommand struct {
	Pat string
}

var _ Command = (*KeysCommand)(nil)

func (k *KeysCommand) Execute(ctx context.Context, env *Env, sess *Session) (parser.Serializable, error) {
	keys := env.Coordinator.Keys(k.Pat)
	res := make(parser.Array[parser.BulkString], len(keys))
	for i, k := range keys {
		res[i] = parser.BulkString(k)
	}
	return res, nil
}

func parseKeys(args []string) (Command, error) {
	return &KeysCommand{Pat: args[0]}, nil
}

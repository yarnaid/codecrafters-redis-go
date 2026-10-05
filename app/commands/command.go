package commands

import (
	"context"

	"my-redis/app/parser"
)

type Command interface {
	Execute(ctx context.Context, env *Env, sess *Session) (parser.Serializable, error)
}

package commands

import (
	"context"
	"fmt"

	"my-redis/app/parser"
)

type ExecCommand struct {
	InTransaction bool
}

var _ Command = (*ExecCommand)(nil)

func (e ExecCommand) Execute(_ context.Context, _ *Env, sess *Session) (parser.Serializable, error) {
	if sess.InMulti {
		return parser.Array[parser.Serializable]{}, nil
	}
	return nil, fmt.Errorf("ERR EXEC without MULTI")
}

func parseExec(args []string) (Command, error) {
	return ExecCommand{}, nil
}

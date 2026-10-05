package commands

import (
	"context"
	"strconv"
	"strings"

	"my-redis/app/parser"
)

type InfoCommand struct{}

var _ Command = (*InfoCommand)(nil)

func (s InfoCommand) Execute(_ context.Context, env *Env, sess *Session) (parser.Serializable, error) {
	data := []string{
		"role:" + string(env.Repl.Role()),
		"master_replid:" + string(env.Repl.ReplID()),
		"master_repl_offset:" + strconv.FormatInt(int64(env.Repl.Offset()), 10),
	}
	return parser.BulkString(strings.Join(data, "")), nil
}

func parseInfo(args []string) (Command, error) {
	if strings.ToLower(args[0]) == "replication" {
		return InfoCommand{}, nil
	}
	return nil, &InvalidCommandError{args[0], args[1:]}
}

package commands

import (
	"context"
	"strconv"
	"strings"

	"my-redis/app/parser"
)

type ReplconfCommand struct {
	Subcmd string
}

var _ Command = (*ReplconfCommand)(nil)

func (r ReplconfCommand) Execute(ctx context.Context, env *Env, sess *Session) (parser.Serializable, error) {
	c := strings.ToLower(r.Subcmd)
	switch c {
	case "listening-port":
	case "capa":
	case "getack":
		return parser.CommandFromStrings("REPLCONF", "ACK", strconv.FormatInt(int64(env.Repl.Offset()), 10)), nil
	}
	return parser.SimpleString("OK"), nil
}

func (r ReplconfCommand) formatResponse(data []string) parser.Serializable {
	res := make(parser.Array[parser.Serializable], len(data)+1)
	res[0] = parser.BulkString("REPLCONF")
	for i, v := range data {
		res[i+1] = parser.BulkString(v)
	}
	return res
}

func parseReplconf(args []string) (Command, error) {
	return ReplconfCommand{Subcmd: args[0]}, nil
}

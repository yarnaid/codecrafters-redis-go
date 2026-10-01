package server

import (
	"strings"

	"my-redis/app/commands"
	"my-redis/app/parser"
)

type ReplconfCommand struct {
	State  *ServerState
	Subcmd string
}

var _ commands.Command = (*ReplconfCommand)(nil)

func (r *ReplconfCommand) Execute() (parser.Serializable, error) {
	c := strings.ToLower(r.Subcmd)
	switch c {
	case "listening-port":
	case "capa":
	}
	return parser.SimpleString("OK"), nil
}

func (r *ReplconfCommand) formatResponse(data []string) parser.Serializable {
	res := make(parser.Array[parser.Serializable], len(data)+1)
	res[0] = parser.BulkString("REPLCONF")
	for i, v := range data {
		res[i+1] = parser.BulkString(v)
	}
	return res
}

func (r *ReplconfCommand) Validate() error {
	return nil
}

package server

import (
	"fmt"

	"my-redis/app/commands"
	"my-redis/app/parser"
)

type PSyncCommand struct {
	ReplId serverId
}

var _ commands.Command = (*PSyncCommand)(nil)

func (p *PSyncCommand) Execute() (parser.Serializable, error) {
	return parser.SimpleString(fmt.Sprintf("FULLRESYNC %s 0", p.ReplId)), nil
}

func (p *PSyncCommand) Validate() error { return nil }

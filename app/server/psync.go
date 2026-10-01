package server

import (
	"my-redis/app/commands"
	"my-redis/app/parser"
)

type PSyncCommand struct{}

var _ commands.Command = (*PSyncCommand)(nil)

func (p *PSyncCommand) Execute() (parser.Serializable, error) {
	return parser.SimpleString("FULLRESYNC <REPL_ID> 0"), nil
}

func (p *PSyncCommand) Validate() error { return nil }

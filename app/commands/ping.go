package commands

import (
	"github.com/codecrafters-io/redis-starter-go/app/parser"
)

type PingCommand struct{}

func (p PingCommand) Execute() (parser.Serializable, error) {
	return parser.SimpleString("PONG"), nil
}
func (p PingCommand) Validate() error {
	return nil
}

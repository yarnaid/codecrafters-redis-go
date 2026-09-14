package commands

import (
	"errors"

	"github.com/codecrafters-io/redis-starter-go/app/parser"
)

type EchoCommand struct {
	Msg string
}

func (e EchoCommand) Execute() (parser.Serializable, error) {
	return parser.BulkString(e.Msg), nil
}

func (e EchoCommand) Validate() error {
	if e.Msg == "" {
		return errors.New("ECHO command requires a message argument")
	}
	return nil
}

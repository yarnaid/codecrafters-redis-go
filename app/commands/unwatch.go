package commands

import (
	"my-redis/app/parser"
)

type UnwatchCommand struct{}

func (w *UnwatchCommand) Execute() (parser.Serializable, error) {
	return parser.SimpleString("OK"), nil
}

func (w *UnwatchCommand) Validate() error {
	return nil
}

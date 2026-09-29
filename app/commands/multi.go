package commands

import (
	"my-redis/app/parser"
)

type MultiCommand struct{}

func (e MultiCommand) Execute() (parser.Serializable, error) {
	return parser.SimpleString("OK"), nil
}

func (e MultiCommand) Validate() error {
	return nil
}

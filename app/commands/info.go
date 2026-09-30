package commands

import (
	"my-redis/app/parser"
)

type InfoCommand struct{}

func (s *InfoCommand) Execute() (parser.Serializable, error) {
	return nil, nil
}

func (s *InfoCommand) Validate() error {
	return nil
}

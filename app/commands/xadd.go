package commands

import (
	"my-redis/app/parser"
	"my-redis/app/storage"
)

type XaddCommand struct {
	Coord  *storage.Coordinator
	Id     string
	Values []parser.Serializable
}

func (s *XaddCommand) Execute() (parser.Serializable, error) {
	s.Coord.Xadd(s.Id, s.Values...)
	return parser.BulkString("0-1"), nil
}

func (s *XaddCommand) Validate() error {
	return nil
}

package commands

import (
	"fmt"

	"my-redis/app/parser"
	"my-redis/app/storage"
)

type SetCommand struct {
	Storage *storage.Storage
	Key     string
	Value   interface{}
	TTL_MS  int
}

func (s *SetCommand) Execute() (parser.Serializable, error) {
	ok := s.Storage.Set(s.Key, s.Value, s.TTL_MS)
	if !ok {
		return nil, fmt.Errorf("Cannot set `%v` to key `%v`", s.Value, s.Key)
	}
	return parser.SimpleString("OK"), nil
}

func (s *SetCommand) Validate() error {
	return nil
}

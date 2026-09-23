package commands

import (
	"fmt"
	"time"

	"my-redis/app/parser"
	"my-redis/app/storage"
)

type SetCommand struct {
	Storage *storage.Coordinator
	Key     string
	Value   interface{}
	TTL     time.Duration
}

func (s *SetCommand) Execute() (parser.Serializable, error) {
	ok := s.Storage.Set(s.Key, s.Value, s.TTL, storage.KindString)
	if !ok {
		return nil, fmt.Errorf("Cannot set `%v` to key `%v`", s.Value, s.Key)
	}
	return parser.SimpleString("OK"), nil
}

func (s *SetCommand) Validate() error {
	return nil
}

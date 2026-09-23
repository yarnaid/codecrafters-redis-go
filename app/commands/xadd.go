package commands

import (
	"my-redis/app/parser"
	"my-redis/app/storage"
)

type XaddCommand struct {
	Coord  *storage.Coordinator
	Key    string
	Id     string
	Values []parser.Serializable
}

func (s *XaddCommand) Execute() (parser.Serializable, error) {
	newKey, err := s.Coord.Xadd(s.Key, s.Id, s.Values...)
	if err != nil {
		switch err.(type) {
		case storage.InvalidStreamIdSeq:
			return parser.SimpleError("ERR The ID specified in XADD is equal or smaller than the target stream top item"), nil
		case storage.ZeroStreamId:
			return parser.SimpleError("ERR The ID specified in XADD must be greater than 0-0"), nil
		default:
			return nil, err
		}
	}
	return parser.BulkString(newKey), nil
}

func (s *XaddCommand) Validate() error {
	return nil
}

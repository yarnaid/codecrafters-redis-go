package commands

import (
	"my-redis/app/parser"
	"my-redis/app/storage"
)

type LLen struct {
	S   *storage.Coordinator
	Key string
}

func (l LLen) Execute() (parser.Serializable, error) {
	return parser.Int(l.S.LLen(l.Key)), nil
}

func (l LLen) Validate() error {
	return nil
}

package commands

import (
	"my-redis/app/parser"
	"my-redis/app/storage"
)

type LPop struct {
	S   *storage.Storage
	Key string
}

func (l LPop) Execute() (parser.Serializable, error) {
	v, err := l.S.LPop(l.Key)
	vv := parser.ToSerializable(v)
	return vv, err
}

func (l LPop) Validate() error {
	return nil
}

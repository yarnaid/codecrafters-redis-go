package commands

import (
	"my-redis/app/parser"
	"my-redis/app/storage"
)

type LPop struct {
	S   *storage.Coordinator
	Key string
	N   int
}

func (l LPop) Execute() (parser.Serializable, error) {
	v, err := l.S.LPop(l.Key, l.N)
	var vv parser.Serializable
	if len(v) == 1 {
		vv = parser.ToSerializable(v[0])
	} else {
		vv = parser.ToSerializable(v)
	}
	return vv, err
}

func (l LPop) Validate() error {
	return nil
}

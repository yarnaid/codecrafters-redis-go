package commands

import (
	"my-redis/app/parser"
	"my-redis/app/storage"
)

type TypeCommand struct {
	Coord *storage.Coordinator
	Key   string
}

func (t *TypeCommand) Execute() (parser.Serializable, error) {
	k, err := t.Coord.Type(t.Key)
	if err != nil {
		return nil, err
	}
	s := string(k)
	return parser.SimpleString(s), nil
}

func (t *TypeCommand) Validate() error {
	return nil
}

package commands

import (
	"my-redis/app/parser"
	"my-redis/app/storage"
)

type IncrCommand struct {
	Coord *storage.Coordinator
	Key   string
}

func (i *IncrCommand) Execute() (parser.Serializable, error) {
	v, err := i.Coord.Incr(i.Key)
	if err != nil {
		switch err.(type) {
		case *storage.WrongTypeError:
			return parser.SimpleError("ERR value is not an integer or out of range"), nil
		default:
			return nil, err
		}
	}
	return v, nil
}

func (i *IncrCommand) Validate() error {
	return nil
}

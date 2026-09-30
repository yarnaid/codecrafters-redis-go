package commands

import (
	"my-redis/app/parser"
	"my-redis/app/storage"
)

type Watch struct {
	Coordinator *storage.Coordinator
	Keys        []string
}

func (w *Watch) Execute() (parser.Serializable, error) {
	return parser.SimpleString("OK"), nil
}

func (w *Watch) Validate() error {
	return nil
}

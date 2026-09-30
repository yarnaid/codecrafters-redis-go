package commands

import (
	"fmt"

	"my-redis/app/parser"
	"my-redis/app/storage"
)

type Watch struct {
	Coordinator *storage.Coordinator
	Keys        []string
	Versions    []int
	InTx        bool
}

func (w *Watch) Execute() (parser.Serializable, error) {
	if w.InTx {
		return nil, fmt.Errorf("ERR WATCH inside MULTI is not allowed")
	}
	w.Versions = w.Coordinator.GetVersions(w.Keys...)
	return parser.SimpleString("OK"), nil
}

func (w *Watch) Validate() error {
	return nil
}

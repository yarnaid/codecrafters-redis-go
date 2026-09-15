package commands

import (
	"my-redis/app/parser"
	"my-redis/app/storage"
)

type RPushCommand struct {
	S      *storage.Storage
	Key    string
	Values []interface{}
}

func (r *RPushCommand) Execute() (parser.Serializable, error) {
	return parser.Int(r.S.Append(r.Key, r.Values...)), nil
}

func (r *RPushCommand) Validate() error {
	return nil
}

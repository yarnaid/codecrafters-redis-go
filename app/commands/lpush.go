package commands

import (
	"my-redis/app/parser"
	"my-redis/app/storage"
)

type LPushCommand struct {
	S      *storage.Coordinator
	Key    string
	Values []parser.Serializable
}

func (r *LPushCommand) Execute() (parser.Serializable, error) {
	return parser.Int(r.S.Prepend(r.Key, r.Values...)), nil
}

func (r *LPushCommand) Validate() error {
	return nil
}

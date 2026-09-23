package commands

import (
	"log/slog"

	"my-redis/app/parser"
	"my-redis/app/storage"
)

type RPushCommand struct {
	S      *storage.Coordinator
	Key    string
	Values []parser.Serializable
}

func (r *RPushCommand) Execute() (parser.Serializable, error) {
	slog.Debug("[Command][RPUSH]", "key", r.Key, "values", r.Values)
	return parser.Int(r.S.Append(r.Key, r.Values...)), nil
}

func (r *RPushCommand) Validate() error {
	return nil
}

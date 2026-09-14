package commands

import (
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/app/parser"
	"github.com/codecrafters-io/redis-starter-go/app/storage"
)

type NotFoundError struct {
	Key string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("Key %v not found in the storage", e.Key)
}

type GetCommand struct {
	strg *storage.Storage
	key  string
}

func (p *GetCommand) Execute() (parser.Serializable, error) {
	val, ok := p.strg.Get(p.key)
	if !ok {
		return parser.BulkNullString(""), nil
	}
	return parser.BulkString(val.(string)), nil
}
func (p *GetCommand) Validate() error {
	return nil
}

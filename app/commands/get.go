package commands

import (
	"fmt"

	"my-redis/app/parser"
	"my-redis/app/storage"
)

type NotFoundError struct {
	Key string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("Key %v not found in the storage", e.Key)
}

type GetCommand struct {
	Strg *storage.Coordinator
	Key  string
}

func (c *GetCommand) Execute() (parser.Serializable, error) {
	// slog.Debug("[Command][GET]", "key", c.Key)
	val, ok := c.Strg.Get(c.Key)
	if !ok {
		// slog.Debug("[Command][GET] value is not received", "key", c.Key)
		return parser.BulkNullString(""), nil
	}
	// slog.Debug("[Command][GET] got", "val", val, "ok", ok)
	switch val := val.(type) {
	case string:
		// slog.Debug("[Command][GET]", "string", val)
		return parser.BulkString(val), nil
	case int:
		// slog.Debug("[Command][GET]", "int", val)
		return parser.Int(val), nil
	default:
		return nil, fmt.Errorf("unsupported type %T", val)
	}
}

func (c *GetCommand) Validate() error {
	return nil
}

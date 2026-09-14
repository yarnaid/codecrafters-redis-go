package commands

import (
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/app/parser"
	"github.com/codecrafters-io/redis-starter-go/app/storage"
)

type SetCommand struct {
	strg  *storage.Storage
	key   string
	value interface{}
}

func (p *SetCommand) Execute() (parser.Serializable, error) {
	ok := p.strg.Set(p.key, p.value)
	if !ok {
		return nil, fmt.Errorf("Cannot set `%v` to key `%v`", p.value, p.key)
	}
	return parser.SimpleString("OK"), nil
}
func (p *SetCommand) Validate() error {
	return nil
}

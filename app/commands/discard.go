package commands

import (
	"fmt"

	"my-redis/app/parser"
)

type DiscardCommand struct {
	TxStarted bool
}

func (e DiscardCommand) Execute() (parser.Serializable, error) {
	if e.TxStarted {
		return parser.SimpleString("OK"), nil
	}
	return nil, fmt.Errorf("ERR DISCARD without MULTI")
}

func (e DiscardCommand) Validate() error {
	return nil
}

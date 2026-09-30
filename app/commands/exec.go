package commands

import (
	"fmt"

	"my-redis/app/parser"
)

type ExecCommand struct {
	InTransaction bool
}

func (e ExecCommand) Execute() (parser.Serializable, error) {
	if e.InTransaction {
		return parser.Array[parser.Serializable]{}, nil
	}
	return nil, fmt.Errorf("ERR EXEC without MULTI")
}

func (e ExecCommand) Validate() error {
	return nil
}

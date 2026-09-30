package commands

import (
	"my-redis/app/parser"
)

type ExecCommand struct {
	InTransaction bool
}

func (e ExecCommand) Execute() (parser.Serializable, error) {
	if e.InTransaction {
		return parser.Array[parser.Serializable]{}, nil
	}
	return parser.SimpleError("ERR EXEC without MULTI"), nil
}

func (e ExecCommand) Validate() error {
	return nil
}

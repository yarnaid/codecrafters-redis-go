package commands

import (
	"my-redis/app/parser"
)

type Command interface {
	Execute() (parser.Serializable, error)
	Validate() error
}

package commands

import (
	"github.com/codecrafters-io/redis-starter-go/app/parser"
)

type Command interface {
	Execute() (parser.Serializable, error)
	Validate() error
}

package commands

import (
	"context"
	"strings"

	"my-redis/app/parser"
)

type ACLCommand struct {
	SubCmd string
	Args   []string
}

var _ Command = (*ACLCommand)(nil)

func (a *ACLCommand) Execute(ctx context.Context, env *Env, sess *Session) (parser.Serializable, error) {
	switch a.SubCmd {
	case "WHOAMI":
		return parser.BulkString("default"), nil
	case "GETUSER":
		return parser.Array[parser.Serializable]{parser.BulkString("flags"), parser.Array[parser.BulkString]{}}, nil
	}
	return nil, &InvalidCommandError{}
}

func parseACL(args []string) (Command, error) {
	switch strings.ToUpper(args[0]) {
	case "WHOAMI":
	case "GETUSER":
	default:
		return nil, &InvalidCommandError{"ACL", args}
	}
	return &ACLCommand{SubCmd: args[0], Args: args[1:]}, nil
}

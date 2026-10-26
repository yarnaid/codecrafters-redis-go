package commands

import (
	"context"
	"fmt"
	"strings"

	"my-redis/app/parser"
	"my-redis/app/user"
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
		passwords := []string{}
		if !sess.User.Flags.Has(user.FlagNoPass) {
			passwords = append(passwords, user.Default().PasswdHash)
		}
		return parser.Array[parser.Serializable]{
			parser.BulkString("flags"), parser.StringsToArray(sess.User.Flags.List()...),
			parser.BulkString("passwords"), parser.StringsToArray(passwords...),
		}, nil
	case "SETUSER":
		if a.Args[1][0] != '>' {
			return nil, fmt.Errorf("incorrect password format, must start from `>`")
		}
		sess.User.SetPassword(a.Args[1][1:])
		return parser.SimpleString("OK"), nil
	}
	return nil, &InvalidCommandError{}
}

func parseACL(args []string) (Command, error) {
	subcmd := strings.ToUpper(args[0])
	switch subcmd {
	case "WHOAMI":
	case "GETUSER":
	case "SETUSER":
	default:
		return nil, &InvalidCommandError{"ACL", args}
	}
	return &ACLCommand{SubCmd: subcmd, Args: args[1:]}, nil
}

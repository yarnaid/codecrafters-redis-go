package commands

import (
	"context"
	"fmt"

	"my-redis/app/parser"
)

type AuthCommand struct {
	Username string
	Password string
}

var _ Command = (*AuthCommand)(nil)

func (a *AuthCommand) Execute(ctx context.Context, env *Env, sess *Session) (parser.Serializable, error) {
	u := env.Users[a.Username]
	if u.CheckPassword(a.Password) {
		sess.User = u
		return parser.SimpleString("OK"), nil
	}
	return nil, fmt.Errorf("WRONGPASS invalid username-password pair or user is disabled.")
}

func parseAuth(args []string) (Command, error) {
	return &AuthCommand{args[0], args[1]}, nil
}

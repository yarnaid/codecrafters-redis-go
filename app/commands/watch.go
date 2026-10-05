package commands

import (
	"context"

	"my-redis/app/parser"
)

type WatchCommand struct {
	Keys []string
}

var _ Command = (*WatchCommand)(nil)

func (w WatchCommand) Execute(ctx context.Context, env *Env, sess *Session) (parser.Serializable, error) {
	versions := env.Coordinator.GetVersions(w.Keys...)
	for i := range versions {
		sess.WatchList = append(sess.WatchList, WatchItem{w.Keys[i], versions[i]})
	}
	return parser.SimpleString("OK"), nil
}

func parseWatch(args []string) (Command, error) {
	return WatchCommand{Keys: args[0:]}, nil
}

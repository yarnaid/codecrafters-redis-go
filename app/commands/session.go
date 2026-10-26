package commands

import (
	"context"

	"my-redis/app/client"
	"my-redis/app/replication"
	"my-redis/app/user"
)

type QueueItem struct {
	Cmd  Command
	Spec Spec
	Args []string
}

type Mode int

const (
	ModeNormal Mode = iota
	ModeSub
)

type Session struct {
	InMulti   bool
	Queue     []QueueItem
	Role      replication.Role
	WatchList []WatchItem
	Client    *client.Client
	Register  func()
	Mode      Mode
	User      *user.User

	Detached bool                                                // stop replica conn handling, not close
	Watch    func(ctx context.Context) (context.Context, func()) // cancel ctx on client disconnect
}

func NewSession(role replication.Role, c *client.Client, u *user.User) *Session {
	if u == nil && user.Default().Flags.Has(user.FlagNoPass) {
		u = user.Default()
	}
	return &Session{
		Queue:     make([]QueueItem, 0),
		Role:      role,
		WatchList: make([]WatchItem, 0),
		Client:    c,
		User:      u,
	}
}

type WatchItem struct {
	Key     string
	Version int
}

func (s *Session) FinishTx() {
	s.InMulti = false
	clear(s.Queue)
	s.Queue = s.Queue[:0]
	s.WatchList = s.WatchList[:0]
}

package commands

import (
	"context"

	"my-redis/app/replication"
)

type QueueItem struct {
	Cmd  Command
	Spec Spec
	Args []string
}

type Session struct {
	InMulti   bool
	Queue     []QueueItem
	Role      replication.Role
	WatchList []WatchItem
	Register  func()

	Detached bool                                                // stop replica conn handling, not close
	Watch    func(ctx context.Context) (context.Context, func()) // cancel ctx on client disconnect
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

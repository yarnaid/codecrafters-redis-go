package commands

import (
	"my-redis/app/parser"
	"my-redis/app/storage"
)

type LRangeCommand struct {
	S          *storage.Coordinator
	Key        string
	Start, End int
}

func (l *LRangeCommand) Execute() (parser.Serializable, error) {
	return l.S.LRange(l.Key, l.Start, l.End), nil
}

func (l *LRangeCommand) Validate() error {
	return nil
}

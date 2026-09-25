package commands

import (
	"my-redis/app/parser"
	"my-redis/app/storage"
)

type XRangeCommand struct {
	C          *storage.Coordinator
	Key        string
	Start, End string
}

func (l *XRangeCommand) Execute() (parser.Serializable, error) {
	slice, err := l.C.XRange(l.Key, l.Start, l.End)
	if err != nil {
		return nil, err
	}
	res := make(parser.Array[parser.Serializable], len(slice))
	for i, v := range slice {
		res[i] = v
	}
	return res, nil
}

func (l *XRangeCommand) Validate() error {
	return nil
}

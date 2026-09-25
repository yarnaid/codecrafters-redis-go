package commands

import (
	"my-redis/app/parser"
	"my-redis/app/storage"
)

type XReadStreamCommand struct {
	C   *storage.Coordinator
	Key string
	Id  string
}

func (l *XReadStreamCommand) Execute() (parser.Serializable, error) {
	slice, err := l.C.XReadStreams(l.Key, l.Id)
	if err != nil {
		return nil, err
	}
	resArr := make(parser.Array[parser.Serializable], len(slice))
	for i, v := range slice {
		resArr[i] = v
	}
	resItem := parser.Array[parser.Serializable]{parser.BulkString(l.Key), resArr}
	res := parser.Array[parser.Serializable]{resItem}
	return res, nil
}

func (l *XReadStreamCommand) Validate() error {
	return nil
}

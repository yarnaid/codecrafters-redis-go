package commands

import (
	"my-redis/app/parser"
	"my-redis/app/storage"
)

type XReadStreamCommand struct {
	C         *storage.Coordinator
	KeyAndIds []string
	ln        int
}

func (l *XReadStreamCommand) Execute() (parser.Serializable, error) {
	l.ln = int(len(l.KeyAndIds) / 2)
	res := make(parser.Array[parser.Serializable], l.ln)
	for i := range l.ln {
		r, err := l.processKey(i)
		if err != nil {
			return nil, err
		}
		res[i] = r
	}
	return res, nil
}

func (l *XReadStreamCommand) processKey(i int) (parser.Serializable, error) {
	k := l.KeyAndIds[i]
	id := l.KeyAndIds[l.ln+i]
	slice, err := l.C.XReadStreams(k, id)
	if err != nil {
		return nil, err
	}
	resArr := make(parser.Array[parser.Serializable], len(slice))
	for i, v := range slice {
		resArr[i] = v
	}
	resItem := parser.Array[parser.Serializable]{parser.BulkString(k), resArr}
	return resItem, nil
}

func (l *XReadStreamCommand) Validate() error {
	return nil
}

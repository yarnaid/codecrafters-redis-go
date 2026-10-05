package commands

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"my-redis/app/parser"
	"my-redis/app/storage"
)

type XReadStreamCommand struct {
	KeyAndIds []string
	Timeout   *time.Duration
	ln        int
}

var _ Command = (*XReadStreamCommand)(nil)

func (l XReadStreamCommand) Execute(ctx context.Context, env *Env, _ *Session) (parser.Serializable, error) {
	l.ln = int(len(l.KeyAndIds) / 2)
	res := make(parser.Array[parser.Serializable], l.ln)
	for i := range l.ln {
		r, err := l.processKey(i, env)
		if err != nil {
			switch err.(type) {
			case *storage.TimeoutError:
				res[i] = parser.NullArray(0)
				continue
			default:
				return nil, err
			}
		}
		res[i] = r
	}
	if len(res) == 1 && res[0] == parser.NullArray(0) {
		return res[0], nil
	}
	return res, nil
}

func (l XReadStreamCommand) processKey(i int, env *Env) (parser.Serializable, error) {
	k := l.KeyAndIds[i]
	id := l.KeyAndIds[l.ln+i]
	slice, err := env.Coordinator.XReadStreams(k, id, l.Timeout)
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

func parseXRead(args []string) (Command, error) {
	cmd := XReadStreamCommand{}
	for len(args) > 0 {
		switch strings.ToUpper(args[0]) {
		case "BLOCK":
			duration, err := strconv.Atoi(args[1])
			if err != nil {
				return nil, err
			}
			timeout := time.Duration(duration) * time.Millisecond
			cmd.Timeout = &timeout
			args = args[2:]
		case "STREAMS":
			cmd.KeyAndIds = args[1:]
			args = args[len(args):]
		default:
			return nil, errors.New("Unknown XREAD subcommand: " + args[1])
		}
	}
	return cmd, nil
}

package server

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	. "my-redis/app/commands"
	"my-redis/app/parser"
)

func (h *Handler) FromInput(input []string) (Command, error) {
	if len(input) == 0 {
		return nil, errors.New("empty command")
	}

	switch strings.ToUpper(input[0]) {
	case "PING":
		return &PingCommand{}, nil
	case "ECHO":
		if len(input) < 2 {
			return nil, errors.New("ECHO command requires an argument")
		}
		return &EchoCommand{Msg: input[1]}, nil
	case "GET":
		if len(input) < 2 {
			return nil, errors.New("GET command requires a key")
		}

		return &GetCommand{Coordinator: h.Coordinator, Key: input[1]}, nil
	case "SET":
		switch len(input) {
		case 3:
			return &SetCommand{Storage: h.Coordinator, Key: input[1], Value: input[2]}, nil
		case 5:
			if strings.ToUpper(input[3]) == "PX" {
				ttl_ms, err := strconv.Atoi(input[4])
				if err != nil {
					return nil, fmt.Errorf("cannot candle set command TTL: %w", err)
				}
				return &SetCommand{Storage: h.Coordinator, Key: input[1], Value: input[2], TTL: time.Duration(ttl_ms) * time.Millisecond}, nil
			}
			return nil, fmt.Errorf("SET got wrong arg: %v", input[3])
		default:
			return nil, fmt.Errorf("SET got wrong args number: %d", len(input))
		}

	case "RPUSH":
		if len(input) < 3 {
			return nil, fmt.Errorf("RPUSH requires at least 2 args, %d got", len(input))
		}
		return &RPushCommand{h.Coordinator, input[1], toSerializableSlice(input[2:])}, nil

	case "LPUSH":
		if len(input) < 3 {
			return nil, fmt.Errorf("LPUSH requires at least 2 args, %d got", len(input))
		}
		return &LPushCommand{h.Coordinator, input[1], toSerializableSlice(input[2:])}, nil

	case "LRANGE":
		if len(input) != 4 {
			return nil, fmt.Errorf("LRANGE args format: key, start, end")
		}
		start, err := strconv.Atoi(input[2])
		if err != nil {
			return nil, err
		}
		end, err := strconv.Atoi(input[3])
		if err != nil {
			return nil, err
		}
		return &LRangeCommand{h.Coordinator, input[1], start, end}, nil

	case "LLEN":
		if len(input) != 2 {
			return nil, fmt.Errorf("wrong args number %d, must be 1 only", len(input)-1)
		}
		return &LLen{h.Coordinator, input[1]}, nil
	case "LPOP":
		if len(input) == 2 {
			return &LPop{h.Coordinator, input[1], 1}, nil
		} else if len(input) == 3 {
			n, err := strconv.Atoi(input[2])
			if err != nil {
				return nil, err
			}
			return &LPop{h.Coordinator, input[1], n}, nil
		}
		return nil, fmt.Errorf("wrong args number %d, must be 1 only", len(input)-1)

	case "BLPOP":
		timeout, err := strconv.ParseFloat(input[2], 64)
		if err != nil {
			return nil, err
		}
		return &BLPop{S: h.Coordinator, Key: input[1], Timeout: time.Duration(math.Round(timeout * float64(time.Second)))}, nil

	case "TYPE":
		return &TypeCommand{Coord: h.Coordinator, Key: input[1]}, nil

	case "XADD":
		return &XaddCommand{Coord: h.Coordinator, Key: input[1], Id: input[2], Values: input[3:]}, nil

	case "XRANGE":
		return &XRangeCommand{C: h.Coordinator, Key: input[1], Start: input[2], End: input[3]}, nil

	case "XREAD":
		input = input[1:]
		cmd := &XReadStreamCommand{C: h.Coordinator}
		for len(input) > 0 {
			switch strings.ToUpper(input[0]) {
			case "BLOCK":
				duration, err := strconv.Atoi(input[1])
				if err != nil {
					return nil, err
				}
				timeout := time.Duration(duration) * time.Millisecond
				cmd.Timeout = &timeout
				input = input[2:]
			case "STREAMS":
				cmd.KeyAndIds = input[1:]
				input = input[len(input):]
			default:
				return nil, errors.New("Unknown XREAD subcommand: " + input[1])
			}
		}
		return cmd, nil

	case "MULTI":
		return &MultiCommand{}, nil

	case "EXEC":
		return &ExecCommand{InTransaction: h.TransactionStarted}, nil

	case "DISCARD":
		return &DiscardCommand{TxStarted: h.TransactionStarted}, nil

	case "INCR":
		return &IncrCommand{Coord: h.Coordinator, Key: input[1]}, nil

	case "WATCH":
		return &WatchCommand{Coordinator: h.Coordinator, Keys: input[1:], InTx: h.TransactionStarted}, nil

	case "UNWATCH":
		return &UnwatchCommand{}, nil

	case "INFO":
		if strings.ToLower(input[1]) == "replication" {
			return &InfoCommand{}, nil
		}
		return &UnwatchCommand{}, &InvalidCommandError{input[0], input[1:]}

	case "REPLCONF":
		return &ReplconfCommand{Subcmd: input[1], State: h.serverState}, nil

	case "PSYNC":
		return &PSyncCommand{ReplId: h.serverState.id}, nil

	case "COMMAND":
		return &InfoCommand{}, nil

	default:
		return nil, &UnknownCommandError{input[0]}
	}
}

func ToAny[T any](s []T) []any {
	result := make([]any, len(s))
	for i, v := range s {
		result[i] = v
	}
	return result
}

func toSerializableSlice[T any](vals []T) []parser.Serializable {
	res := make([]parser.Serializable, len(vals))
	for i, val := range vals {
		res[i] = parser.ToSerializable(val)
	}
	return res
}

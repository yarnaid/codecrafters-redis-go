package commands

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"my-redis/app/parser"
	"my-redis/app/storage"
)

func FromInput(input []string, coordinator *storage.Coordinator) (Command, error) {
	if len(input) == 0 {
		return nil, errors.New("Empty command")
	}

	switch input[0] {
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

		return &GetCommand{Strg: coordinator, Key: input[1]}, nil
	case "SET":
		switch len(input) {
		case 3:
			return &SetCommand{Storage: coordinator, Key: input[1], Value: input[2]}, nil
		case 5:
			if strings.ToUpper(input[3]) == "PX" {
				ttl_ms, err := strconv.Atoi(input[4])
				if err != nil {
					return nil, fmt.Errorf("cannot candle set command TTL: %w", err)
				}
				return &SetCommand{Storage: coordinator, Key: input[1], Value: input[2], TTL: time.Duration(ttl_ms) * time.Millisecond}, nil
			}
			return nil, fmt.Errorf("SET got wrong arg: %v", input[3])
		default:
			return nil, fmt.Errorf("SET got wrong args number: %d", len(input))
		}

	case "RPUSH":
		if len(input) < 3 {
			return nil, fmt.Errorf("RPUSH requires at least 2 args, %d got", len(input))
		}
		return &RPushCommand{coordinator, input[1], toSerializableSlice(input[2:])}, nil

	case "LPUSH":
		if len(input) < 3 {
			return nil, fmt.Errorf("LPUSH requires at least 2 args, %d got", len(input))
		}
		return &LPushCommand{coordinator, input[1], toSerializableSlice(input[2:])}, nil

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
		return &LRangeCommand{coordinator, input[1], start, end}, nil

	case "LLEN":
		if len(input) != 2 {
			return nil, fmt.Errorf("wrong args number %d, must be 1 only", len(input)-1)
		}
		return &LLen{coordinator, input[1]}, nil
	case "LPOP":
		if len(input) == 2 {
			return &LPop{coordinator, input[1], 1}, nil
		} else if len(input) == 3 {
			n, err := strconv.Atoi(input[2])
			if err != nil {
				return nil, err
			}
			return &LPop{coordinator, input[1], n}, nil
		}
		return nil, fmt.Errorf("wrong args number %d, must be 1 only", len(input)-1)

	case "BLPOP":
		timeout, err := strconv.Atoi(input[2])
		if err != nil {
			return nil, err
		}
		return &BLPop{S: coordinator, Key: input[1], Timeout: time.Duration(timeout) * time.Second}, nil

	default:
		return nil, errors.New("Unknown command: " + input[0])
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

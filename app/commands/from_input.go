package commands

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"my-redis/app/storage"
)

func FromInput(input []string, strg *storage.Storage) (Command, error) {
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

		return &GetCommand{Strg: strg, Key: input[1]}, nil
	case "SET":
		switch len(input) {
		case 3:
			return &SetCommand{Storage: strg, Key: input[1], Value: input[2]}, nil
		case 5:
			if strings.ToUpper(input[3]) == "PX" {
				ttl_ms, err := strconv.Atoi(input[4])
				if err != nil {
					return nil, fmt.Errorf("cannot candle set command TTL: %w", err)
				}
				return &SetCommand{Storage: strg, Key: input[1], Value: input[2], TTL_MS: ttl_ms}, nil
			}
			return nil, fmt.Errorf("SET got wrong arg: %v", input[3])
		default:
			return nil, fmt.Errorf("SET got wrong args number: %d", len(input))
		}

	case "RPUSH":
		if len(input) < 3 {
			return nil, fmt.Errorf("RPUSH requires at least 2 args, %d got", len(input))
		}
		return &RPushCommand{strg, input[1], ToAny(input[2:])}, nil

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
		return &LRangeCommand{strg, input[1], start, end}, nil
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

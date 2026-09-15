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

	default:
		return nil, errors.New("Unknown command: " + input[0])
	}
}

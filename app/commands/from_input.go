package commands

import (
	"errors"

	"github.com/codecrafters-io/redis-starter-go/app/storage"
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
		return &GetCommand{strg: strg, key: input[1]}, nil
	case "SET":
		if len(input) != 3 {
			return nil, errors.New("SET requires 2 args: key and value")
		}
		return &SetCommand{strg: strg, key: input[1], value: input[2]}, nil
	default:
		return nil, errors.New("Unknown command: " + input[0])
	}
}

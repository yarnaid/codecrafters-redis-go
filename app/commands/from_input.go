package commands

import (
	"errors"
)

func FromInput(input []string) (Command, error) {
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
	default:
		return nil, errors.New("Unknown command: " + input[0])
	}
}

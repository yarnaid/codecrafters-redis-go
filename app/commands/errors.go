package commands

import "fmt"

type NotFoundError struct {
	Key string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("Key %v not found in the storage", e.Key)
}

type InvalidCommandError struct {
	Cmd  string
	Args []string
}

func (i *InvalidCommandError) Error() string {
	return fmt.Sprintf("invalid command `%s` with args %v", i.Cmd, i.Args)
}

type UnknownCommandError struct {
	Cmd string
}

func (i *UnknownCommandError) Error() string {
	return "unknown command " + i.Cmd
}

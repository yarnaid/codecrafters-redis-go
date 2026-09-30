package commands

import "fmt"

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

package commands

import "fmt"

type Flags uint8

const (
	FlagWrite Flags = 1 << iota
	FlagBlocking
	FlagExclusive
	FlagNotPropagate
	FlagAllowedInSubs
)

func (f Flags) Propagate() bool {
	return f&FlagNotPropagate == 0
}

func (f Flags) Has(prop Flags) bool {
	return f&prop == prop
}

// Spec describes one command. Parse builds a fresh Command per request.
type Spec struct {
	Name    string
	MinArgs int // excluding the command name
	MaxArgs int
	Flags   Flags
	Parse   func(args []string) (Command, error)
}

func (s Spec) Validate(args []string) error {
	if len(args) < s.MinArgs || (s.MaxArgs > 0 && len(args) > s.MaxArgs) {
		return fmt.Errorf("%v: incorrect args number %d,", s.Name, len(args))
	}
	return nil
}

func (s Spec) ParseArgs(args []string) (Command, error) {
	if err := s.Validate(args); err != nil {
		return nil, err
	}
	return s.Parse(args)
}

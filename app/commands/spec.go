//go:generate go tool stringer -type=Flags -linecomment
package commands

import (
	"fmt"
	"math/bits"
	"strings"
)

type Flags uint8

const (
	FlagWrite         Flags = 1 << iota // write
	FlagBlocking                        // blocking
	FlagExclusive                       // exclusive
	FlagNotPropagate                    // not_propagate
	FlagAllowedInSubs                   // allowed_in_subs
)

func (f Flags) Propagate() bool {
	return f&FlagNotPropagate == 0
}

func (f Flags) Has(prop Flags) bool {
	return f&prop == prop
}

func FlagsToString(f Flags) string {
	if f == 0 {
		return "none"
	}
	var parts []string
	// Iterate set bits from lowest to highest: deterministic order.
	for v := uint8(f); v != 0; v &= v - 1 {
		bit := Flags(1 << bits.TrailingZeros8(v)) // math/bits
		parts = append(parts, fmt.Sprint(bit))
	}
	return strings.Join(parts, "|")
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

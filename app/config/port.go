package config

import (
	"fmt"
	"strconv"
)

type Port uint16

func (p Port) MarshalText() ([]byte, error) {
	return strconv.AppendUint(nil, uint64(p), 10), nil
}

func (p *Port) UnmarshalText(b []byte) error {
	v, err := strconv.ParseUint(string(b), 10, 16) // bitSize 16 enforces 0-65535
	if err != nil {
		return fmt.Errorf("invalid port %q: want 0-65535", b)
	}
	*p = Port(v)
	return nil
}

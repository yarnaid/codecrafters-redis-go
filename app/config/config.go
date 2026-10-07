// Package config ...
package config

import (
	"net/netip"
)

type Config struct {
	Bind       netip.Addr
	Port       Port
	Replicaof  string
	Debug      bool
	Dir        string
	DBFilename string
}

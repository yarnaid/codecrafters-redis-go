package main

import (
	"math/rand/v2"
	"strconv"
	"strings"

	"my-redis/app/parser"
)

type serverState struct {
	role   role
	offset int
	id     serverId
}

func (s *serverState) info() []string {
	return []string{
		"role:" + string(s.role),
		"master_replid:" + string(s.id),
		"master_repl_offset:" + strconv.FormatInt(int64(s.offset), 10),
	}
}

func (s *serverState) toBulkString() parser.BulkString {
	str := strings.Join(s.info(), "\n")
	return parser.BulkString(str)
}

type role string

const (
	masterRole role = "master"
	slaveRole  role = "slave"
)

type serverId string

func newServerId() serverId {
	l := 40
	chars := make([]string, l)
	for i := range l {
		chars[i] = string(rune('a' + rand.IntN(10)))
	}
	return serverId(strings.Join(chars, ""))
}

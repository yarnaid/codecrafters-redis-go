package replication

import (
	"math/rand/v2"
	"strings"
)

type ReplID string

func newReplID() ReplID {
	l := 40
	chars := make([]string, l)
	for i := range l {
		chars[i] = string(rune('a' + rand.IntN('z'-'a'+1)))
	}
	return ReplID(strings.Join(chars, ""))
}

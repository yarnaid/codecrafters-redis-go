package replication

import (
	"math/rand/v2"
	"strings"
)

type replID string

func newReplID() replID {
	l := 40
	chars := make([]string, l)
	for i := range l {
		chars[i] = string(rune('a' + rand.IntN('z'-'a'+1)))
	}
	return replID(strings.Join(chars, ""))
}

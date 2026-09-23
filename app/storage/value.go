package storage

import (
	"time"
)

type Value struct {
	ExpireAt time.Time
	Value    interface{}
	Kind     ValueKind
}

func (v *Value) IsExpired() bool {
	if !v.ExpireAt.IsZero() && v.ExpireAt.Before(time.Now()) {
		return true
	}
	return false
}

package storage

import (
	"time"
)

type Value struct {
	ExpireAt time.Time
	Value    interface{}
}

func (v *Value) IsExpired() bool {
	if !v.ExpireAt.IsZero() && v.ExpireAt.Before(time.Now()) {
		return true
	}
	return false
}

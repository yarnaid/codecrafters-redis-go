package storage

import (
	"fmt"
	"time"
)

type ValueKind string

const (
	KindString ValueKind = "string"
	KindList   ValueKind = "list"
	KindNone   ValueKind = "none"
)

func (k ValueKind) IsValid() bool {
	switch k {
	case KindList, KindString:
		return true
	}
	return false
}

func ParseValueKind(s string) (ValueKind, error) {
	k := ValueKind(s)
	if !k.IsValid() {
		return "", fmt.Errorf("invalid ValueKind %q", s)
	}
	return k, nil
}

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

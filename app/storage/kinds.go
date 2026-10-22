package storage

import "fmt"

type ValueKind string

const (
	KindString ValueKind = "string"
	KindList   ValueKind = "list"
	KindStream ValueKind = "stream"
	KindNone   ValueKind = "none"
	KindInt    ValueKind = "int"
	KindBitMap ValueKind = "bitmap"
)

func (k ValueKind) IsValid() bool {
	switch k {
	case KindList, KindString, KindStream, KindNone, KindBitMap:
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

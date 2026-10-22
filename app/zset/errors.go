package zset

import "fmt"

type NotFoundError struct {
	key string
}

var _ error = (*NotFoundError)(nil)

func (n *NotFoundError) Error() string {
	return fmt.Sprintf("key not found `%s`", n.key)
}

func notFoundErrorf(key string) error {
	return &NotFoundError{key: key}
}

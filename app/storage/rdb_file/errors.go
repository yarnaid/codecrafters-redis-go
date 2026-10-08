package rdbfile

import (
	"errors"
	"fmt"
)

type ParseError struct {
	Msg string
	Err error
}

var _ error = (*ParseError)(nil)

func (e *ParseError) Error() string {
	res := "rdb parse: " + e.Msg
	if e.Err != nil {
		res += ": " + e.Err.Error()
	}
	return res
}

func (e *ParseError) Unwrap() error {
	return e.Err
}

func parseErrorf(cause error, msg string, args ...any) *ParseError {
	return &ParseError{Msg: fmt.Sprintf(msg, args...), Err: cause}
}

var (
	ErrLZF           = errors.New("LZF-compressed strings are not supported")
	ErrUnknownOpCode = errors.New("unknown opcode")
)

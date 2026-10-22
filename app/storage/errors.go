package storage

import (
	"fmt"
	"time"
)

type TimeoutError struct {
	Key      string
	At       time.Time
	Duration time.Duration
	StreamID *StreamID
}

func (t *TimeoutError) Error() string {
	return fmt.Sprintf("timeout error: key=%v, at=%v, duration=%v, stream_id=%v", t.Key, t.At, t.Duration, deref(t.StreamID))
}

func NewTimeoutError(key string, duration time.Duration, streamID *StreamID) *TimeoutError {
	return &TimeoutError{
		Key:      key,
		At:       time.Now(),
		Duration: duration,
		StreamID: streamID,
	}
}

type WrongTypeError struct {
	Got      any
	Required any
}

func (w *WrongTypeError) Error() string {
	return fmt.Sprintf("wrong type: expected=%T, got=%T", w.Required, w.Got)
}

type NotFoundError struct {
	Key string
}

func (n *NotFoundError) Error() string {
	return fmt.Sprintf("key=%v not found", n.Key)
}

func (n *NotFoundError) Is(target error) bool {
	_, ok := target.(*NotFoundError)
	return ok
}

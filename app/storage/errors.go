package storage

import (
	"fmt"
	"time"
)

type TimeoutError struct {
	Key      string
	At       time.Time
	Duration time.Duration
	StreamId *StreamId
}

func (t *TimeoutError) Error() string {
	return fmt.Sprintf("timeout error: key=%v, at=%v, duration=%v, stream_id=%v", t.Key, t.At, t.Duration, deref(t.StreamId))
}

func NewTimeoutError(key string, duration time.Duration, streamId *StreamId) *TimeoutError {
	return &TimeoutError{
		Key:      key,
		At:       time.Now(),
		Duration: duration,
		StreamId: streamId,
	}
}

type WrongTypeError struct {
	Got      interface{}
	Required interface{}
}

func (w *WrongTypeError) Error() string {
	return fmt.Sprintf("wrong type: expected=%T, got=%T", w.Required, w.Got)
}

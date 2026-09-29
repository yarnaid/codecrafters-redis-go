package storage

import (
	"container/list"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"my-redis/app/parser"
)

func (s *Coordinator) Append(key string, values ...parser.Serializable) int {
	slog.Debug("[Coordinator][Append]", "key", key, "values", values)
	s.mu.Lock()
	defer s.mu.Unlock()
	var popped int

	if handoff := s.tryHandover(key, values[0]); handoff {
		slog.Debug("[Coordinator][Append] handoff success for value 0", "key", key, "values", values)
		values = values[1:]
		popped = 1
	}

	val, ok := s.backend.Get(key)
	if !ok {
		slog.Debug("[Coordinator][Append] no array, creating", "key", key)
		// slog.Debug("[Storage][Append] current data", "data", s.data)
		val = Value{Value: make([]parser.Serializable, 0), Kind: KindList}
		s.backend.Set(key, &val)
	}
	arr, ok := val.Value.([]parser.Serializable)
	if !ok {
		return 0
	}
	arr = append(arr, values...)
	val.Value = arr
	s.backend.Set(key, &val)
	slog.Debug("[Coordinator][Append] len", "len", len(arr), "val", val.Value)
	return len(arr) + popped
}

func (s *Coordinator) Prepend(key string, values ...parser.Serializable) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	var popped int

	if handoff := s.tryHandover(key, values[0]); handoff {
		slog.Debug("[Coordinator][Prepend] handoff success, removing value 0", "key", key, "values", values)
		values = values[1:]
		popped = 1
	}

	val, ok := s.backend.Get(key)
	if !ok {
		// slog.Debug("[Coordinator][Append] no array, creating", "key", key)
		val = Value{Value: make([]parser.Serializable, 0), Kind: KindList}
		s.backend.Set(key, &val)
	}
	arr, ok := val.Value.([]parser.Serializable)
	if !ok {
		return 0
	}
	val.Value = PrependReversed(arr, values...)
	s.backend.Set(key, &val)
	res_arr, _ := val.Value.([]parser.Serializable)
	// slog.Debug("[Storage][Append] len", "len", len(res_arr), "val", res_arr)
	return len(res_arr) + popped
}

func (s *Coordinator) LRange(key string, start, end int) parser.Array[parser.Serializable] {
	s.mu.RLock()
	defer s.mu.RUnlock()

	val, ok := s.backend.Get(key)
	if !ok {
		slog.Debug("[Coordinator][LRange] key not found")
		return parser.Array[parser.Serializable]{}
	}
	arr, ok := val.Value.([]parser.Serializable)
	if !ok {
		slog.Debug("[Coordinator][LRange] key is not array", "type", fmt.Sprintf("%T", val.Value))
		return parser.Array[parser.Serializable]{}
	}
	res := rangeList(arr, start, end)
	slog.Debug("[Coordinator][LRange] returning", "res", res, "start", start, "end", end, "orig", arr)
	return res
}

func (s *Coordinator) LLen(key string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	val, ok := s.backend.Get(key)
	if !ok {
		slog.Debug("[Coordinator][LLen] key not found")
		return 0
	}
	arr, ok := val.Value.([]parser.Serializable)
	if !ok {
		slog.Debug("[Coordinator][LLen] key is not array")
		return 0
	}
	return len(arr)
}

func (s *Coordinator) LPop(key string, n int) ([]parser.Serializable, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	val, ok := s.backend.Get(key)
	if !ok {
		slog.Debug("[Coordinator][LPop] key not found")
		return nil, errors.New("not found")
	}
	arr, ok := val.Value.([]parser.Serializable)
	if !ok {
		slog.Debug("[Coordinator][LRange] key is not array")
		return nil, errors.New("not an array")
	}
	if len(arr) == 0 {
		return nil, errors.New("arr is empty")
	}

	slog.Debug("[Coordinator][LPop] input", "arr", arr)
	res, remain := popListLeft(arr, n)
	s.backend.Set(key, &Value{Value: remain})
	slog.Debug("[Coordinator][LPop] return", "res", res, "new_arr", remain)
	return res, nil
}

func (c *Coordinator) BLPop(key string, timeout time.Duration) (parser.Serializable, error) {
	slog.Debug("[Coordinator][BLPOP]", "key", key, "timeout", timeout)
	c.mu.Lock()
	raw_val, ok := c.backend.Get(key)
	if ok {
		slog.Debug("[Coordinator][BLPOP] simple way", "key", key, "timeout", timeout)

		list, ok := raw_val.Value.([]parser.Serializable)
		if !ok {
			c.mu.Unlock()
			return nil, fmt.Errorf("cannot convert %v to serializable slice (%T)", raw_val.Value, raw_val.Value)
		}
		if len(list) > 0 {
			popped, remained := popListLeft(list, 1)
			if len(remained) > 0 {
				c.backend.Set(key, &Value{Value: remained})
			} else {
				c.backend.Delete(key)
			}
			c.mu.Unlock()
			return popped[0], nil
		}
	}
	return c.waitForResult(key, timeout)
}

func (c *Coordinator) waitForResult(key string, timeout time.Duration) (parser.Serializable, error) {
	slog.Debug("[Coordinator][waitForResult] blocking way", "key", key, "timeout", timeout)

	var timeout_ch <-chan time.Time
	if timeout > 0 {
		timer := time.NewTimer(timeout)
		defer timer.Stop()
		timeout_ch = timer.C
	}
	result_waiter := NewWaiterArray()
	waiters, ok := c.blpopWaiters[key]
	if !ok {
		waiters = list.New()
		c.blpopWaiters[key] = waiters
	}
	waiters.PushBack(result_waiter)
	result_waiter.elem = waiters.Back()
	c.blpopWaiters[key] = waiters
	c.mu.Unlock()

	slog.Debug("[Coordinator][waitForResult] start waiting", "key", key, "waiters", waiters.Len())
	for {
		select {
		case <-timeout_ch:
			slog.Debug("[Coordinator][waitForResult] timer fired")
			c.mu.Lock()
			slog.Debug("[Coordinator][waitForResult] timer under lock")
			waiters.Remove(result_waiter.elem)
			if waiters.Len() == 0 {
				delete(c.blpopWaiters, key)
			}
			c.mu.Unlock()
			return nil, NewTimeoutError(key, timeout, nil)
		case val := <-result_waiter.ch:
			return val, nil
		}
	}
}

func (c *Coordinator) tryHandover(key string, value parser.Serializable) bool {
	waiters, ok := c.blpopWaiters[key]
	if !ok || waiters.Len() == 0 {
		slog.Debug("[Coordinator][tryHandoff] no waiters found", "key", key, "waiters", c.blpopWaiters)
		return false
	}

	front := waiters.Front()
	waiter := front.Value.(waiterArray)
	waiters.Remove(front)
	if waiters.Len() == 0 {
		delete(c.blpopWaiters, key)
	}
	waiter.ch <- value
	slog.Debug("[Coordinator][tryHandoff] success", "key", key)
	return true
}

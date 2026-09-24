package storage

import (
	"container/list"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"my-redis/app/parser"
)

type Coordinator struct {
	mu            sync.RWMutex
	backend       Backend
	blpop_waiters map[string]*list.List
}

type waiter struct {
	ch   chan parser.Serializable
	elem *list.Element
}

func NewWaiter() waiter {
	return waiter{ch: make(chan parser.Serializable, 1)}
}

func NewMemoryCoordinator(memory *MemoryBackend) *Coordinator {
	var mem *MemoryBackend
	if memory != nil {
		mem = memory
	} else {
		mem = &MemoryBackend{data: make(map[string]*Value)}
	}
	return &Coordinator{
		backend:       mem,
		blpop_waiters: make(map[string]*list.List, 0),
	}
}

func (s *Coordinator) Set(key string, value interface{}, ttl time.Duration, kind ValueKind) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	// slog.Debug("[Coordinator][Set]", "key", key, "value", value, "ttl", ttl_ms)
	val := Value{Value: value, Kind: kind}
	if ttl > 0 {
		val.ExpireAt = time.Now().Add(ttl)
	}
	s.backend.Set(key, &val)
	return true
}

func (s *Coordinator) Get(key string) (interface{}, bool) {
	// slog.Debug("[Coordinator][GET]", "key", key)
	s.mu.RLock()
	defer s.mu.RUnlock()

	val, ok := s.backend.Get(key)
	// slog.Debug("[Coordinator][GET] got", "key", key, "val", val, "ok", ok)
	if !ok {
		return nil, ok
	}
	if val.IsExpired() {
		slog.Debug("[Coordinator][GET] expired", "ago", time.Since(val.ExpireAt))
		s.backend.Delete(key)
		return nil, false
	}
	return val.Value, true
}

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
	slog.Debug("[Coordinator][BLPOP] blocking way", "key", key, "timeout", timeout)

	var timeout_ch <-chan time.Time
	if timeout > 0 {
		timer := time.NewTimer(timeout)
		defer timer.Stop()
		timeout_ch = timer.C
	}
	result_waiter := NewWaiter()
	waiters, ok := c.blpop_waiters[key]
	if !ok {
		waiters = list.New()
		c.blpop_waiters[key] = waiters
	}
	waiters.PushBack(result_waiter)
	result_waiter.elem = waiters.Back()
	c.blpop_waiters[key] = waiters
	c.mu.Unlock()

	slog.Debug("[Coordinator][BLPOP] start waiting", "key", key, "waiters", waiters.Len())
	select {
	case <-timeout_ch:
		slog.Debug("[Coordinator][BLPOP] timer fired")
		c.mu.Lock()
		slog.Debug("[Coordinator][BLPOP] timer under lock")
		waiters.Remove(result_waiter.elem)
		if waiters.Len() == 0 {
			delete(c.blpop_waiters, key)
		}
		c.mu.Unlock()
		return nil, fmt.Errorf("blpop timeout=%v, key=%s", timeout, key)
		// res := parser.ToSerializable(<-result_waiter.ch)
		// slog.Debug("[Coordinator][BLPOP] returning blocked value")
		// return res, nil
	case val := <-result_waiter.ch:
		res := parser.ToSerializable(val)
		return res, nil
	}
}

func (c *Coordinator) tryHandover(key string, value parser.Serializable) bool {
	waiters, ok := c.blpop_waiters[key]
	if !ok || waiters.Len() == 0 {
		slog.Debug("[Coordinator][tryHandoff] no waiters found", "key", key, "waiters", c.blpop_waiters)
		return false
	}

	front := waiters.Front()
	waiter := front.Value.(waiter)
	waiters.Remove(front)
	if waiters.Len() == 0 {
		delete(c.blpop_waiters, key)
	}
	waiter.ch <- value
	slog.Debug("[Coordinator][tryHandoff] success", "key", key)
	return true
}

func (c *Coordinator) Type(key string) (ValueKind, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, ok := c.backend.Get(key)
	if !ok {
		return KindNone, nil
	}
	if !v.Kind.IsValid() {
		slog.Error("incorrect value kind", "kind", v.Kind)
		return "", fmt.Errorf("invalid kind %q", v.Kind)
	}
	return v.Kind, nil
}

func (c *Coordinator) Xadd(key, id string, kvPairs ...parser.Serializable) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	v, ok := c.backend.Get(key)
	if !ok {
		v = Value{Kind: KindStream, Value: make([]*StreamContainer, 0)}
	}
	stream, ok := v.Value.([]*StreamContainer)
	if !ok {
		return id, fmt.Errorf("value %T is not slice stream containers k=%q, v=%v", v.Value, key, v.Value)
	}
	newId, err := ParseStreamId(id)
	if err != nil {
		return id, err
	}
	newStreamVal := &StreamContainer{Id: newId}
	if len(stream) > 0 {
		last := stream[len(stream)-1]
		if last.GreaterOrGen(newStreamVal) || (last.Id.Seq == newStreamVal.Id.Seq && last.Id.Time == newStreamVal.Id.Time) {
			err = InvalidStreamIdSeq{&last.Id, &newStreamVal.Id}
			slog.Error("invalid stream id seq", "err", err)
			return newStreamVal.Id.String(), err
		}
		slog.Debug("[Coordinator][Xadd] greater!", "key", key, "oldId", last.Id, "newId", newId)
	} else {
		slog.Debug("[Coordinator][Xadd] skipping empty list validation", "key", key)
		if newStreamVal.Id.Time < 0 {
			newStreamVal.Id = StreamId{int(time.Now().UnixMilli()), 0}
		} else {
			if newStreamVal.Id.Seq < 0 {
				if newStreamVal.Id.Time == 0 {
					newStreamVal.Id.Seq = 1
				} else {
					newStreamVal.Id.Seq = 0
				}
			}
		}
	}
	v.Value = append(stream, newStreamVal)
	c.backend.Set(key, &v)
	return newStreamVal.Id.String(), nil
}

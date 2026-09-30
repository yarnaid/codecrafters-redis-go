package storage

import (
	"container/list"
	"fmt"
	"log/slog"
	"strconv"
	"sync"
	"time"

	"my-redis/app/parser"
)

type Coordinator struct {
	mu           sync.RWMutex
	backend      Backend
	blpopWaiters map[string]*list.List
	xreadWaiters map[string]*list.List
}

func NewMemoryCoordinator(memory *MemoryBackend) *Coordinator {
	var mem *MemoryBackend
	if memory != nil {
		mem = memory
	} else {
		mem = &MemoryBackend{data: make(map[string]*Value)}
	}
	return &Coordinator{
		backend:      mem,
		blpopWaiters: make(map[string]*list.List, 0),
		xreadWaiters: make(map[string]*list.List, 0),
	}
}

func (s *Coordinator) Set(key string, value interface{}, ttl time.Duration, kind ValueKind) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.setLocked(key, value, ttl, kind)
}

func (s *Coordinator) setLocked(key string, value interface{}, ttl time.Duration, kind ValueKind) bool {
	version := s.getVersionsLocked(key)[0]
	val := Value{Value: value, Kind: kind, Version: max(version+1, 1)}
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

	val, _, err := s.getLocked(key)
	return val, err
}

func (s *Coordinator) getLocked(key string) (interface{}, *time.Time, bool) {
	val, ok := s.backend.Get(key)
	if !ok {
		return nil, nil, ok
	}
	if val.IsExpired() {
		slog.Debug("[Coordinator][GET] expired", "ago", time.Since(val.ExpireAt))
		s.backend.Delete(key)
		return nil, nil, false
	}
	return val.Value, &val.ExpireAt, true
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

func (c *Coordinator) Incr(key string) (parser.Int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	val, expire_at, ok := c.getLocked(key)
	if !ok {
		val = int(0)
	}
	var val_int int
	switch v := val.(type) {
	case string:
		var err error
		val_int, err = strconv.Atoi(v)
		if err != nil {
			switch err.(type) {
			case *strconv.NumError:
				return 0, &WrongTypeError{Got: "", Required: int(0)}
			default:
				return 0, err
			}
		}
	case int:
		val_int = v
	default:
		return 0, &WrongTypeError{Got: val, Required: ""}
	}
	val_int++
	var ttl time.Duration
	if expire_at != nil {
		ttl = time.Until(*expire_at)
	} else {
		ttl = 0
	}
	ok = c.setLocked(key, val_int, ttl, KindInt)
	if !ok {
		return 0, fmt.Errorf("cannot set value=%d for key=%v", val_int, key)
	}
	return parser.Int(val_int), nil
}

func (c *Coordinator) GetVersions(keys ...string) []int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.getVersionsLocked(keys...)
}

func (c *Coordinator) getVersionsLocked(keys ...string) []int {
	versions := make([]int, len(keys))
	for i, k := range keys {
		v, ok := c.backend.Get(k)
		if !ok {
			versions[i] = -1
		} else {
			versions[i] = v.Version
		}
	}
	return versions
}

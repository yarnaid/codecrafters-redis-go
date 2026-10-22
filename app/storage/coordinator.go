// Package storage ...
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

func (c *Coordinator) Set(key string, value interface{}, ttl time.Duration, kind ValueKind) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.setLocked(key, value, ttl, kind)
}

func (c *Coordinator) setLocked(key string, value interface{}, ttl time.Duration, kind ValueKind) bool {
	version := c.getVersionsLocked(key)[0]
	val := Value{Value: value, Kind: kind, Version: max(version+1, 1)}
	if ttl > 0 {
		val.ExpireAt = time.Now().Add(ttl)
	}
	c.backend.Set(key, &val)
	return true
}

func (c *Coordinator) Get(key string) (interface{}, bool) {
	// slog.Debug("[Coordinator][GET]", "key", key)
	c.mu.RLock()
	defer c.mu.RUnlock()

	val, _, err := c.getLocked(key)
	return val, err
}

func (c *Coordinator) getLocked(key string) (interface{}, *time.Time, bool) {
	val, ok := c.backend.Get(key)
	if !ok {
		return nil, nil, ok
	}
	if val.IsExpired() {
		slog.Debug("[Coordinator][GET] expired", "ago", time.Since(val.ExpireAt))
		c.backend.Delete(key)
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
	val, expireAt, ok := c.getLocked(key)
	if !ok {
		val = int(0)
	}
	var valInt int
	switch v := val.(type) {
	case string:
		var err error
		valInt, err = strconv.Atoi(v)
		if err != nil {
			switch err.(type) {
			case *strconv.NumError:
				return 0, &WrongTypeError{Got: "", Required: int(0)}
			default:
				return 0, err
			}
		}
	case int:
		valInt = v
	default:
		return 0, &WrongTypeError{Got: val, Required: ""}
	}
	valInt++
	var ttl time.Duration
	if expireAt != nil {
		ttl = time.Until(*expireAt)
	} else {
		ttl = 0
	}
	ok = c.setLocked(key, valInt, ttl, KindInt)
	if !ok {
		return 0, fmt.Errorf("cannot set value=%d for key=%v", valInt, key)
	}
	return parser.Int(valInt), nil
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

func (c *Coordinator) Keys(pat string) []string {
	return c.backend.Keys(pat)
}

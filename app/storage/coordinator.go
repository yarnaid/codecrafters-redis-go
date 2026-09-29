package storage

import (
	"container/list"
	"fmt"
	"log/slog"
	"sync"
	"time"
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

// deref returns *p, or the string "<nil>" for logging purposes.
func deref[T any](p *T) any {
	if p == nil {
		return "<nil>"
	}
	return *p
}

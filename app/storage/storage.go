package storage

import (
	"log/slog"
	"sync"
	"time"
)

var _ = slog.Log

type Storage struct {
	mu   sync.RWMutex
	data map[string]*Value
}

func NewStorage() *Storage {
	res := Storage{
		data: make(map[string]*Value),
	}
	return &res
}

func (s *Storage) Set(key string, value interface{}, ttl_ms int) bool {
	s.mu.Lock()
	defer func() {
		s.mu.Unlock()
	}()
	slog.Debug("storage set", "key", key, "value", value, "ttl", ttl_ms)
	val := Value{Value: value}
	if ttl_ms > 0 {
		val.ExpireAt = time.Now().Add(time.Duration(ttl_ms) * time.Millisecond)
	}
	s.data[key] = &val
	return true
}

func (s *Storage) Get(key string) (interface{}, bool) {
	slog.Debug("[storage] GET", "key", key)
	s.mu.RLock()
	defer func() {
		s.mu.RUnlock()
	}()
	val, ok := s.data[key]
	slog.Debug("[storage] got", "key", key, "val", val, "ok", ok)
	if !ok {
		return nil, ok
	}
	if val.IsExpired() {
		slog.Debug("[Storage][GET] expired", "ago", time.Since(val.ExpireAt))
		delete(s.data, key)
		return nil, false
	}
	return val.Value, true
}

func (s *Storage) Append(key string, values ...interface{}) int {
	s.mu.Lock()
	defer func() {
		s.mu.Unlock()
	}()

	val, ok := s.data[key]
	if !ok {
		slog.Debug("[Storage][Append] no array, creating", "key", key)
		// slog.Debug("[Storage][Append] current data", "data", s.data)
		val = &Value{Value: make([]interface{}, 0)}
		s.data[key] = val
	}
	arr, ok := val.Value.([]interface{})
	if !ok {
		return 0
	}
	val.Value = append(arr, values...)
	res_arr, _ := val.Value.([]interface{})
	slog.Debug("[Storage][Append] len", "len", len(res_arr), "val", res_arr)
	return len(res_arr)
}

package storage

import (
	"log/slog"
	"sync"
)

type Storage struct {
	mu   sync.RWMutex
	data map[string]interface{}
}

func NewStorage() *Storage {
	res := Storage{
		data: make(map[string]interface{}),
	}
	return &res
}

func (s *Storage) Set(key string, value interface{}) bool {
	s.mu.Lock()
	defer func() {
		s.mu.Unlock()
	}()
	slog.Debug("storage set", "key", key, "value", value)
	s.data[key] = value
	return true
}

func (s *Storage) Get(key string) (interface{}, bool) {
	s.mu.RLock()
	defer func() {
		s.mu.RUnlock()
	}()
	val, ok := s.data[key]
	if !ok {
		return nil, ok
	}
	return val, true
}

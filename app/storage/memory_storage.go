package storage

import (
	"errors"
	"log/slog"
	"sync"
	"time"

	"my-redis/app/parser"
)

var _ = slog.Log

type MemoryStorage struct {
	mu   sync.RWMutex
	data map[string]*Value
}

func NewMemoryStorage() *MemoryStorage {
	panic("this memory storage is deprecated!")
}

func (s *MemoryStorage) Set(key string, value any, ttlMs int) bool {
	s.mu.Lock()
	defer func() {
		s.mu.Unlock()
	}()
	slog.Debug("storage set", "key", key, "value", value, "ttl", ttlMs)
	val := Value{Value: value}
	if ttlMs > 0 {
		val.ExpireAt = time.Now().Add(time.Duration(ttlMs) * time.Millisecond)
	}
	s.data[key] = &val
	return true
}

func (s *MemoryStorage) Get(key string) (interface{}, bool) {
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

func (s *MemoryStorage) Append(key string, values ...any) int {
	s.mu.Lock()
	defer func() {
		s.mu.Unlock()
	}()

	val, ok := s.data[key]
	if !ok {
		slog.Debug("[Storage][Append] no array, creating", "key", key)
		// slog.Debug("[Storage][Append] current data", "data", s.data)
		val = &Value{Value: make([]any, 0)}
		s.data[key] = val
	}
	arr, ok := val.Value.([]any)
	if !ok {
		return 0
	}
	val.Value = append(arr, values...)
	resArr, _ := val.Value.([]any)
	// slog.Debug("[Storage][Append] len", "len", len(res_arr), "val", res_arr)
	return len(resArr)
}

func (s *MemoryStorage) Prepend(key string, values ...any) int {
	s.mu.Lock()
	defer func() {
		s.mu.Unlock()
	}()

	val, ok := s.data[key]
	if !ok {
		slog.Debug("[Storage][Append] no array, creating", "key", key)
		// slog.Debug("[Storage][Append] current data", "data", s.data)
		val = &Value{Value: make([]any, 0)}
		s.data[key] = val
	}
	arr, ok := val.Value.([]any)
	if !ok {
		return 0
	}
	val.Value = PrependReversed(arr, values...)
	resArr, _ := val.Value.([]any)
	// slog.Debug("[Storage][Append] len", "len", len(res_arr), "val", res_arr)
	return len(resArr)
}

func (s *MemoryStorage) LRange(key string, start, end int) parser.Array[parser.Serializable] {
	s.mu.RLock()
	defer func() {
		s.mu.RUnlock()
	}()

	val, ok := s.data[key]
	if !ok {
		slog.Debug("[Storage][LRange] key not found")
		return parser.Array[parser.Serializable]{}
	}
	arr, ok := val.Value.([]interface{})
	if !ok {
		slog.Debug("[Storage][LRange] key is not array")
		return parser.Array[parser.Serializable]{}
	}
	var res parser.Array[parser.Serializable]
	if end >= 0 {
		end = min(len(arr)-1, end)
	}
	start, end = invertIndex(start, len(arr)), invertIndex(end, len(arr))
	for _, v := range arr[start : end+1] {
		res = append(res, parser.ToSerializable(v))
	}
	// slog.Debug("[Storage][LRange] returning", "res", res, "start", start, "end", end, "orig", arr)
	return res
}

func (s *MemoryStorage) LLen(key string) int {
	s.mu.RLock()
	defer func() {
		s.mu.RUnlock()
	}()

	val, ok := s.data[key]
	if !ok {
		slog.Debug("[Storage][LLen] key not found")
		return 0
	}
	arr, ok := val.Value.([]interface{})
	if !ok {
		slog.Debug("[Storage][LRange] key is not array")
		return 0
	}
	return len(arr)
}

func (s *MemoryStorage) LPop(key string, n int) ([]interface{}, error) {
	s.mu.RLock()
	defer func() {
		s.mu.RUnlock()
	}()

	val, ok := s.data[key]
	if !ok {
		slog.Debug("[Storage][LPop] key not found")
		return nil, errors.New("not found")
	}
	arr, ok := val.Value.([]interface{})
	if !ok {
		slog.Debug("[Storage][LRange] key is not array")
		return nil, errors.New("not an array")
	}
	if len(arr) == 0 {
		return nil, errors.New("arr is empty")
	}

	slog.Debug("[Storage][LPop] input", "arr", arr)
	nn := min(n, len(arr))
	v := make([]interface{}, nn)
	copy(v, arr[:nn])
	copy(arr[:len(arr)-nn], arr[nn:])
	arr = arr[:len(arr)-nn]
	val.Value = arr
	slog.Debug("[Storage][LPop] return", "res", v, "new_arr", arr)
	return v, nil
}

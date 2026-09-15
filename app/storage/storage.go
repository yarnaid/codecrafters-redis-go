package storage

import (
	"errors"
	"log/slog"
	"sync"
	"time"

	"my-redis/app/parser"
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
	// slog.Debug("[Storage][Append] len", "len", len(res_arr), "val", res_arr)
	return len(res_arr)
}

func (s *Storage) Prepend(key string, values ...interface{}) int {
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
	val.Value = PrependReversed(arr, values...)
	res_arr, _ := val.Value.([]interface{})
	// slog.Debug("[Storage][Append] len", "len", len(res_arr), "val", res_arr)
	return len(res_arr)
}

func PrependReversed[T any](s []T, values ...T) []T {
	result := make([]T, 0, len(values)+len(s))
	for i := len(values) - 1; i >= 0; i-- {
		result = append(result, values[i])
	}
	return append(result, s...)
}

func (s *Storage) LRange(key string, start, end int) parser.Array[parser.Serializable] {
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
	start, end = invert_index(start, len(arr)), invert_index(end, len(arr))
	for _, v := range arr[start : end+1] {
		res = append(res, parser.ToSerializable(v))
	}
	// slog.Debug("[Storage][LRange] returning", "res", res, "start", start, "end", end, "orig", arr)
	return res
}

func invert_index(i, ln int) int {
	if i >= 0 {
		return i
	} else {
		return max(ln+i, 0)
	}
}

func (s *Storage) LLen(key string) int {
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

func (s *Storage) LPop(key string) (interface{}, error) {
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
	v := arr[0]
	copy(arr[:len(arr)-1], arr[1:])
	arr = arr[:len(arr)-1]
	val.Value = arr
	return v, nil
}

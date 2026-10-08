package storage

import (
	"log/slog"
	"maps"
	"slices"
)

type MemoryBackend struct {
	data map[string]*Value
}

var _ Backend = (*MemoryBackend)(nil)

func NewMemoryBackend(data map[string]*Value) *MemoryBackend {
	var data_set map[string]*Value
	if data != nil {
		data_set = data
	} else {
		data_set = make(map[string]*Value)
	}
	return &MemoryBackend{data: data_set}
}

func (m *MemoryBackend) Get(key string) (Value, bool) {
	// slog.Debug("[MemoryBackend][Get]", "key", key, "data", m.data)
	v, ok := m.data[key]
	if !ok {
		slog.Debug("[MemoryBackend] not found", "key", key)
		return Value{}, false
	}
	return *v, ok
}

func (m *MemoryBackend) Set(key string, v *Value) {
	m.data[key] = v
}

func (m *MemoryBackend) Delete(key string) {
	delete(m.data, key)
}

func (m *MemoryBackend) Keys(pat string) []string {
	it := maps.Keys(m.data)
	return slices.AppendSeq([]string{}, it)
}

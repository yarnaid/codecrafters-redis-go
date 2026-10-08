package rdbfile

import "my-redis/app/storage"

type Snapshot struct {
	Version  int
	Metadata map[string]string
	Data     map[string]*storage.Value
}

func NewSnapshot() *Snapshot {
	return &Snapshot{
		Version:  11,
		Data:     make(map[string]*storage.Value),
		Metadata: make(map[string]string),
	}
}

func (s *Snapshot) ToMemoryBackend() *storage.MemoryBackend {
	return storage.NewMemoryBackend(s.Data)
}

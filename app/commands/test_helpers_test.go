package commands_test

import (
	"time"

	"my-redis/app/parser"
	"my-redis/app/storage"
)

type CoordinatorValue struct {
	Key   string
	Value parser.Serializable
	TTL   time.Duration
	Kind  storage.ValueKind
}

func GetCoordinatorWithData(data []CoordinatorValue) (c *storage.Coordinator) {
	c = storage.NewMemoryCoordinator(nil)
	for _, v := range data {
		c.Set(v.Key, v.Value, v.TTL, v.Kind)
	}
	return c
}

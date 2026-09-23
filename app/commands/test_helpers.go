package commands

import (
	"time"

	"my-redis/app/parser"
	"my-redis/app/storage"
)

type CoordinatorValue struct {
	Key   string
	Value parser.Serializable
	TTL   time.Duration
}

func GetCoordinatorWithData(data []CoordinatorValue) (c *storage.Coordinator) {
	c = storage.NewMemoryCoordinator(nil)
	for _, v := range data {
		c.Set(v.Key, v.Value, v.TTL)
	}
	return c
}

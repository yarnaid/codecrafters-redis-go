package storage

import (
	"slices"

	"my-redis/app/zset"
)

func (c *Coordinator) ZAdd(key, name string, score float64) (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	res := 0
	v, ok := c.backend.Get(key)
	if !ok {
		v = Value{Value: zset.NewZSetMapSlice()}
		c.backend.Set(key, &v)
		res = 1
	}
	z, ok := v.Value.(zset.ZSet)
	if !ok {
		return -1, &WrongTypeError{v.Value, zset.ZSetMapSlice{}}
	}
	if z.Add(name, score) {
		return 1, nil
	}
	return res, nil
}

func (c *Coordinator) ZRank(key, name string) (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	got, ok := c.backend.Get(key)
	if !ok {
		return -1, &NotFoundError{name}
	}
	z, ok := got.Value.(zset.ZSet)
	if !ok {
		return -1, &WrongTypeError{got.Value, zset.ZSetMapSlice{}}
	}
	rank, ok := z.Rank(name)
	if !ok {
		return -1, &NotFoundError{name}
	}
	return rank, nil
}

func (c *Coordinator) ZRange(key string, start, end int) ([]string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	got, ok := c.backend.Get(key)
	if !ok {
		return []string{}, nil
	}
	z, ok := got.Value.(zset.ZSet)
	if !ok {
		return nil, &WrongTypeError{got.Value, zset.ZSetMapSlice{}}
	}
	l := z.Count()
	start, end = invert_index(start, l), invert_index(end, l)
	return slices.Collect(z.Range(start, end+1)), nil
}

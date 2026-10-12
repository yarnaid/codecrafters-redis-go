package storage

import "my-redis/app/zset"

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

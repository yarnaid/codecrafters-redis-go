package storage

import "my-redis/app/storage/bitmap"

func (c *Coordinator) SetBit(key string, offset int, zero bool) (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	val, ok := c.backend.Get(key)
	if !ok {
		val = Value{Value: bitmap.NewBitMap()}
		c.backend.Set(key, &val)
	}

	b, ok := val.Value.(*bitmap.BitMap)
	if !ok {
		return -1, &WrongTypeError{b, &bitmap.BitMap{}}
	}
	old := b.Bit(offset)
	if zero {
		b.Clear(offset)
	} else {
		b.Set(offset)
	}

	return old, nil
}

package storage

import (
	"my-redis/app/storage/bitmap"
)

func (c *Coordinator) GetBit(key string, offset int) (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	val, ok := c.backend.Get(key)
	if !ok {
		return 0, nil
	}

	var b *bitmap.BitMap
	switch v := val.Value.(type) {
	case *bitmap.BitMap:
		b = v
	case string:
		b = bitmap.FromString(v)
	default:
		return -1, &WrongTypeError{b, &bitmap.BitMap{}}
	}

	return b.Bit(offset), nil
}

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

func (c *Coordinator) BitCount(key string, start, end int) (int, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	v, ok := c.backend.Get(key)
	if !ok {
		// slog.Debug("[Coordinator][bitcount] not found", "key", key)
		return 0, false
	}
	switch vv := v.Value.(type) {
	case *bitmap.BitMap:
		return vv.BitCount(start, end), true
	case string:
		return bitmap.FromString(vv).BitCount(start, end), true
	}
	// slog.Debug("[Coordinator] wrong type", "key", key, "v", v.Value, "type", fmt.Sprintf("%t", v.Value))
	return 0, false
}

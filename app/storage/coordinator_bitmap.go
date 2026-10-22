package storage

import (
	"fmt"
	"strings"

	"my-redis/app/storage/bitmap"
)

func (c *Coordinator) GetBit(key string, offset int) (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	b, err := c.getBitMapLocked(key)
	if err != nil {
		return 0, err
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

	b, err := c.getBitMapLocked(key)
	if err != nil {
		return 0, false
	}
	return b.BitCount(start, end), true
}

func (c *Coordinator) getBitMapLocked(key string) (*bitmap.BitMap, error) {
	v, ok := c.backend.Get(key)
	if !ok {
		return nil, fmt.Errorf("not found")
	}
	var b *bitmap.BitMap
	switch vv := v.Value.(type) {
	case *bitmap.BitMap:
		b = vv
	case string:
		b = bitmap.FromString(vv)
	default:
		return nil, fmt.Errorf("invalid type for bitmap")
	}
	return b, nil
}

func (c *Coordinator) BitOp(op, dest string, k1, k2 string) (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	b1, _ := c.getBitMapLocked(k1)
	if b1 == nil {
		b1 = bitmap.NewBitMap()
	}
	b2, _ := c.getBitMapLocked(k2)
	if b2 == nil {
		b2 = bitmap.NewBitMap()
	}
	switch strings.ToUpper(op) {
	case "AND":
		res := bitmap.And(b1, b2)
		c.backend.Set(dest, &Value{Value: res})
		return res.DataLen(), nil
	default:
		return 0, fmt.Errorf("not implemented op for bitmap %s", op)
	}
}

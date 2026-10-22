package bitmap

import "math/bits"

func NewBitMap() *BitMap {
	return &BitMap{
		data: []byte{},
	}
}

func FromString(s string) *BitMap {
	b := make([]byte, len(s))
	for i, r := range s {
		b[i] = bits.Reverse8(uint8(r))
	}
	return &BitMap{
		data: b,
		len:  len(s) * 8,
	}
}

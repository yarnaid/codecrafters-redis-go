package zset

import (
	"cmp"
	"slices"
)

func compareVal(a, b *value) int {
	v := cmp.Compare(a.score, b.score)
	switch v {
	case 0:
		return cmp.Compare(a.name, b.name)
	default:
		return v
	}
}

type ZSetMapSlice struct {
	Map  map[string]float64
	vals []*value
}

var _ ZSet = (*ZSetMapSlice)(nil)

func NewZSetMapSlice() *ZSetMapSlice {
	return &ZSetMapSlice{
		Map:  map[string]float64{},
		vals: []*value{},
	}
}

func (z *ZSetMapSlice) Add(name string, s float64) {
	vv := &value{name, s}
	z.Map[name] = s
	idx, exists := slices.BinarySearchFunc(z.vals, vv, compareVal)
	if !exists {
		z.vals = slices.Insert(z.vals, idx, vv)
	}
}

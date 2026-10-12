package zset

import (
	"cmp"
	"log/slog"
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

func (z *ZSetMapSlice) Add(name string, s float64) bool {
	if v, ok := z.Map[name]; ok && v == s {
		return false
	}
	vv := &value{name, s}
	idx, found := slices.BinarySearchFunc(z.vals, vv, compareVal)
	if !found {
		slog.Debug("found: %v, %v", name, z.vals)
		z.vals = slices.Insert(z.vals, idx, vv)
		z.Map[name] = s
		return true
	}
	return false
}

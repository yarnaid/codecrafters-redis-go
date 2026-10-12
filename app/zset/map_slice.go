package zset

import (
	"cmp"
	"slices"
)

func compareVal(a, b *value) int {
	if c := cmp.Compare(a.score, b.score); c != 0 {
		return c
	}
	return cmp.Compare(a.name, b.name)
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
	old, exists := z.Map[name]
	if exists && old == s {
		return false
	}
	newVal := &value{name: name, score: s}
	if exists {
		i, _ := slices.BinarySearchFunc(z.vals, newVal, compareVal)
		if i < len(z.vals) {
			z.vals = slices.Delete(z.vals, i, i+1)
		} else {
			z.vals = z.vals[:len(z.vals)-1]
		}
	}
	i, _ := slices.BinarySearchFunc(z.vals, newVal, compareVal)
	z.vals = slices.Insert(z.vals, i, newVal)
	z.Map[name] = s
	return !exists
}

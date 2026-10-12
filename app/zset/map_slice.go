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

func (z *ZSetMapSlice) binSearch(v *value) (int, bool) {
	return slices.BinarySearchFunc(z.vals, v, compareVal)
}

func (z *ZSetMapSlice) Add(name string, s float64) bool {
	old, exists := z.Map[name]
	if exists && old == s {
		return false
	}
	newVal := &value{name: name, score: s}
	if exists {
		i, _ := z.binSearch(newVal)
		if i < len(z.vals) {
			z.vals = slices.Delete(z.vals, i, i+1)
		} else {
			z.vals = z.vals[:len(z.vals)-1]
		}
	}
	i, _ := z.binSearch(newVal)
	z.vals = slices.Insert(z.vals, i, newVal)
	z.Map[name] = s
	return !exists
}

func (z *ZSetMapSlice) Rank(name string) (int, bool) {
	score, found := z.Map[name]
	if !found {
		return -1, false
	}
	return z.binSearch(&value{name, score})
}

package zset

import (
	"cmp"
	"iter"
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

func (z *ZSetMapSlice) Count() int {
	return len(z.vals)
}

func (z *ZSetMapSlice) Rank(name string) (int, bool) {
	score, found := z.Map[name]
	if !found {
		return -1, false
	}
	return z.binSearch(&value{name, score})
}

func (z *ZSetMapSlice) Range(start, end int) iter.Seq[string] {
	return func(yield func(string) bool) {
		res := z.vals[start:end]
		for _, v := range res {
			if !yield(v.name) {
				return
			}
		}
	}
}

func (z *ZSetMapSlice) Score(key string) (float64, bool) {
	s, ok := z.Map[key]
	return s, ok
}

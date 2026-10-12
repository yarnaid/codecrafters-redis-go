package zset

import (
	"fmt"
	"iter"
	"maps"
	"slices"
)

// Map lazily applies f to each value of seq.
func Map[T, U any](seq iter.Seq[T], f func(T) U) iter.Seq[U] {
	return func(yield func(U) bool) {
		for v := range seq {
			if !yield(f(v)) {
				return // consumer stopped early (break)
			}
		}
	}
}

func (z *ZSetMapSlice) Keys() []string {
	return slices.Sorted(maps.Keys(z.Map))
}

func (z *ZSetMapSlice) ValKeys() []string {
	return slices.Sorted(Map(slices.Values(z.vals), func(v *value) string { return v.name }))
}

func (z *ZSetMapSlice) SliceVals() []float64 {
	return slices.Collect(Map(slices.Values(z.vals), func(v *value) float64 { return v.score }))
}

func (z *ZSetMapSlice) ValidateMapSlice() error {
	if len(z.Map) != len(z.vals) {
		return fmt.Errorf("map(%d) and vals(%d) have different len", len(z.Map), len(z.vals))
	}
	mapKeys := z.Keys()
	valKeys := z.ValKeys()
	if eq := slices.Equal(mapKeys, valKeys); !eq {
		return fmt.Errorf("map and slice have different keys set")
	}
	if !slices.IsSorted(z.SliceVals()) {
		return fmt.Errorf("values in zset not sorted")
	}
	return nil
}

func HasKeys(z *ZSetMapSlice, keys ...string) (bool, error) {
	if err := z.ValidateMapSlice(); err != nil {
		return false, err
	}
	mapKeys := z.Keys()
	for _, k := range keys {
		if !slices.Contains(mapKeys, k) {
			return false, fmt.Errorf("zset doesn't contain key `%s`", k)
		}
	}
	return true, nil
}

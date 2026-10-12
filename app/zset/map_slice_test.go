package zset_test

import (
	"maps"
	"slices"
	"testing"

	"my-redis/app/zset"

	"github.com/stretchr/testify/require"
)

func TestMapSlice_Add(t *testing.T) {
	tests := []struct {
		name string
		kv   map[string]float64
	}{
		{"1 value", map[string]float64{"a": 1.0}},
		{"2 values", map[string]float64{"a": 1.0, "b": 2.0}},
		{"2 values", map[string]float64{"b": 2.0, "a": 1.0}},
		{"3 values", map[string]float64{"c": 10, "b": 2.0, "a": 1.0}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)
			z := zset.NewZSetMapSlice()
			for k, v := range tt.kv {
				z.Add(k, v)
			}
			require.NoError(z.ValidateMapSlice())
			require.EqualValues(slices.Sorted(maps.Keys(tt.kv)), z.Keys())
			require.EqualValues(slices.Sorted(maps.Values(tt.kv)), z.SliceVals())
		})
	}
}

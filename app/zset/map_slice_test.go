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
		new  []bool
	}{
		{"1 value", map[string]float64{"a": 1.0}, []bool{true}},
		{"2 same value", map[string]float64{"a": 1.0}, []bool{true, false}},
		{"2 values", map[string]float64{"a": 1.0, "b": 2.0}, []bool{true, true}},
		{"2 values", map[string]float64{"b": 2.0, "a": 1.0}, []bool{true, true}},
		{"3 values", map[string]float64{"c": 10, "b": 2.0, "a": 1.0}, []bool{true, true, true}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)
			z := zset.NewZSetMapSlice()
			for i, k := range slices.Collect(maps.Keys(tt.kv)) {
				added := z.Add(k, tt.kv[k])
				require.Equal(tt.new[i], added)

			}
			require.NoError(z.ValidateMapSlice())
			require.EqualValues(slices.Sorted(maps.Keys(tt.kv)), z.Keys())
			require.EqualValues(slices.Sorted(maps.Values(tt.kv)), z.SliceVals())
		})
	}
}

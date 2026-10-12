package zset_test

import (
	"slices"
	"testing"

	"my-redis/app/zset"

	"github.com/stretchr/testify/require"
)

func TestMapSlice_Add(t *testing.T) {
	type kv struct {
		k string
		v float64
	}
	tests := []struct {
		name string
		kv   []kv
		new  []bool
	}{
		{"1 value", []kv{{"a", 1.0}}, []bool{true}},
		{"2 same score", []kv{{"a", 1.0}, {"a", 1.0}}, []bool{true, false}},
		{"2 diff score", []kv{{"a", 1.0}, {"a", 2.0}}, []bool{true, false}},
		{"2 values", []kv{{"a", 1.0}, {"b", 2.0}}, []bool{true, true}},
		{"2 values", []kv{{"b", 2.0}, {"a", 1.0}}, []bool{true, true}},
		{"3 values", []kv{{"c", 10}, {"b", 2.0}, {"a", 1.0}}, []bool{true, true, true}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)
			z := zset.NewZSetMapSlice()
			for i, k := range tt.kv {
				added := z.Add(k.k, k.v)
				require.Equal(tt.new[i], added)
			}
			keys := []string{}
			vals := []float64{}
			for _, k := range tt.kv {
				keys = append(keys, k.k)
				vals = append(vals, k.v)
			}
			keys = slices.Sorted(slices.Values(slices.Compact(keys)))
			vals = slices.Sorted(slices.Values(slices.Compact(vals)))
			require.NoError(z.ValidateMapSlice())
			require.EqualValues(keys, z.Keys())
			// require.EqualValues(vals, z.SliceVals())
			// require.Contains(vals, z.SliceVals())
		})
	}
}

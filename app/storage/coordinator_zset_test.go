package storage_test

import (
	"fmt"
	"maps"
	"slices"
	"testing"

	"my-redis/app/storage"

	"github.com/stretchr/testify/require"
)

func TestCoordinator_ZAdd(t *testing.T) {
	type val struct {
		k string
		v float64
	}
	tests := []struct {
		name  string
		kv    []val
		added []int
	}{
		{"1 new val", []val{{"a", 1}}, []int{1}},
		{"1 new val+dup", []val{{"a", 1}, {"a", 1}}, []int{1, 0}},
		{"2 new val+dup", []val{{"a", 1}, {"b", 2}, {"a", 1}}, []int{1, 1, 0}},
		{"2 new val+upd", []val{{"a", 1}, {"b", 2}, {"a", 2}}, []int{1, 1, 0}},
	}
	key := "zset-key"

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)
			c := storage.NewMemoryCoordinator(nil)
			for i, v := range tt.kv {
				added, err := c.ZAdd(key, v.k, v.v)
				require.NoError(err)
				require.Equal(tt.added[i], added)
			}
		})
	}
}

func TestCoordinator_ZRange(t *testing.T) {
	key := "zrange-key"
	data := map[string]float64{
		"0": 0, "1": 1, "2": 2, "3": 3, "4": 4,
	}
	keys := slices.Collect(maps.Keys(data))
	c := storage.NewMemoryCoordinator(nil)
	for k, v := range data {
		if i, err := c.ZAdd(key, k, v); err != nil || i != 1 {
			t.Fatal("cannot create zset", i, err)
		}
	}
	tests := []struct {
		start, end int
		expected   []string
	}{
		{0, 4, keys},
		{0, 10, keys},
		{12, 10, []string{}},
		{0, 2, keys[:3]},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d -> %d", tt.start, tt.end), func(t *testing.T) {
			require := require.New(t)
			actual, err := c.ZRange(key, tt.start, tt.end)
			require.NoError(err)
			require.EqualValues(tt.expected, actual)
		})
	}
}

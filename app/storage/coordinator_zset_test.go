package storage_test

import (
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

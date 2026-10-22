package bitmap_test

import (
	"testing"

	"my-redis/app/storage/bitmap"

	"github.com/stretchr/testify/require"
)

func newBM(data []byte) *bitmap.BitMap {
	b := bitmap.NewBitMap()
	b.SetData(data)
	return b
}

func TestAnd(t *testing.T) {
	tests := []struct {
		name   string // description of this test case
		d1, d2 []byte
		want   []byte
	}{
		{"1&1=1", []byte{1}, []byte{1}, []byte{1}},
		{"1&2=1", []byte{1}, []byte{2}, []byte{0}},
		{"1&2=1", []byte{1, 0}, []byte{2}, []byte{0, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)

			got := bitmap.And(newBM(tt.d1), newBM(tt.d2))
			require.Equal(tt.want, got.Data())
		})
	}
}

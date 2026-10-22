package bitmap_test

import (
	"testing"

	"my-redis/app/storage/bitmap"

	"github.com/stretchr/testify/require"
)

func Test_bitMap_Clear(t *testing.T) {
	data := []byte{
		0b11111111,
		0b11111111,
	}
	tests := []struct {
		name     string
		offset   int
		expected string
	}{
		{"first to zero", 0, "0111111111111111"},
		{"second to zero", 1, "1011111111111111"},
		{"last to zero", 15, "1111111111111110"},
		{"before last to zero", 14, "1111111111111101"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)
			b := bitmap.NewBitMap()
			b.SetData(data)
			b.SetLen(16)
			b.Clear(tt.offset)
			require.Equal(tt.expected, b.String())
		})
	}
}

func Test_bitMap_String(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		data []byte
		len  int
		want string
	}{
		{"zero", []byte{}, 0, ""},
		{"one zero", []byte{}, 1, "0"},
		{"two zero", []byte{}, 2, "00"},
		{"one 1", []byte{0b1}, 1, "1"},
		{"two 1", []byte{0b11}, 2, "11"},
		{"10", []byte{0b1}, 2, "10"},
		{"01", []byte{0b10}, 2, "01"},
		{"010", []byte{0b10}, 3, "010"},
		{"00000000 1", []byte{0, 0b1}, 9, "000000001"},
		// {"2 bytes", []byte{0b0,0b}, 2, "01"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)
			b := bitmap.NewBitMap()
			b.SetData(tt.data)
			b.SetLen(tt.len)
			got := b.String()
			require.Equal(tt.want, got)
		})
	}
}

func Test_bitMap_Extend(t *testing.T) {
	b := bitmap.NewBitMap()

	t.Run("extend", func(t *testing.T) {
		require := require.New(t)
		require.Len(b.Data(), 0)
		b.ExtendTo(1)
		require.Len(b.Data(), 1)
		b.ExtendTo(1)
		require.Len(b.Data(), 1)
		b.ExtendTo(2)
		require.Len(b.Data(), 2)
	})
}

func Test_bitMap_Set(t *testing.T) {
	tests := []struct {
		name     string
		offset   int
		expected string
	}{
		{"first", 0, "1"},
		{"second", 1, "01"},
		{"10", 10, "00000000001"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)
			b := bitmap.NewBitMap()
			b.Set(tt.offset)
			require.Equal(tt.expected, b.String())
		})
	}
}

func TestBitMap_AsString(t *testing.T) {
	tests := []struct {
		data []byte
		want string
	}{
		{[]byte{0b10}, "@"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			require := require.New(t)
			b := bitmap.NewBitMap()
			b.SetData(tt.data)
			b.SetLen(2)
			got := b.AsString()
			require.Equal(tt.want, got)
		})
	}
}

package storage

import (
	"fmt"
	"log/slog"
	"testing"

	"my-redis/app/parser"

	"github.com/stretchr/testify/require"
)

func TestRangeList(t *testing.T) {
	l := 7
	data := make(parser.Array[parser.Serializable], l)
	for i := range l {
		data[i] = parser.Int(i)
	}
	tests := []struct {
		start, end int
		expected   parser.Array[parser.Serializable]
	}{
		{0, 6, data[0:7]},
		{0, 5, data[:6]},
		{0, 0, data[:1]},
		{0, 1, data[:2]},
		{0, -1, data[:l]},
		{0, -2, data[:l-1]},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d, %d", tt.start, tt.end), func(t *testing.T) {
			require := require.New(t)
			actual := rangeList(data, tt.start, tt.end)
			slog.Debug("m", "data", data, "actual", actual)
			require.EqualValues(tt.expected, actual)
		})
	}
}

func TestInvertIndex(t *testing.T) {
	tests := []struct {
		val, ln, expected int
	}{
		{0, 10, 0},
		{-1, 10, 9},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d from %d", tt.val, tt.ln), func(t *testing.T) {
			require := require.New(t)
			actual := invert_index(tt.val, tt.ln)
			require.Equal(tt.expected, actual)
		})
	}
}

func TestPopList(t *testing.T) {
	l := 7
	list := make([]parser.Serializable, l)
	for i := 0; i < l; i++ {
		list[i] = parser.Int(i)
	}
	tests := []struct {
		n               int
		expected_popped []parser.Serializable
		expected_tail   []parser.Serializable
	}{
		{1, list[:1], list[1:]},
		{2, list[:2], list[2:]},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("pop %d", tt.n), func(t *testing.T) {
			require := require.New(t)
			popped, remaining := popListLeft(list, tt.n)
			require.EqualValues(tt.expected_popped, popped)
			require.EqualValues(tt.expected_tail, remaining)
		})
	}
}

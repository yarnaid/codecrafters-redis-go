package bitmap_test

import (
	"fmt"
	"testing"

	"my-redis/app/storage/bitmap"

	"github.com/stretchr/testify/require"
)

func TestFromString(t *testing.T) {
	tests := []struct {
		from     string
		expected string
	}{
		{"A", "01000001"},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s -> %s", tt.from, tt.expected), func(t *testing.T) {
			require := require.New(t)
			b := bitmap.FromString(tt.from)
			require.Equal(tt.expected, b.String())
		})
	}
}

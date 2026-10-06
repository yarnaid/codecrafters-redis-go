package parser_test

import (
	"bytes"
	"testing"

	"my-redis/app/parser"

	"github.com/stretchr/testify/require"
)

func TestReader(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected [][]string
		wantErr  bool
	}{
		{"simple", parser.StringsToBytes("1", "2"), [][]string{{"1", "2"}}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)
			reader := parser.NewReader(bytes.NewReader(tt.input))
			act, err := reader.ReadArrays()
			require.NoError(err)
			require.EqualValues(tt.expected, act)
		})
	}
}

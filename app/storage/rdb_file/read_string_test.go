package rdbfile_test

import (
	"bufio"
	"bytes"
	"testing"

	rdbfile "my-redis/app/storage/rdb_file"

	"github.com/stretchr/testify/require"
)

func TestReadString(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		str   string
		err   error
	}{
		{"simple", []byte{0x0D, 0x48, 0x65, 0x6C, 0x6C, 0x6F, 0x2C, 0x20, 0x57, 0x6F, 0x72, 0x6C, 0x64, 0x21}, "Hello, World!", nil},
		{"int8", []byte{0xC0, 0x7B}, "123", nil},
		{"int16", []byte{0xC1, 0x39, 0x30}, "12345", nil},
		{"int32", []byte{0xC2, 0x87, 0xD6, 0x12, 0x00}, "1234567", nil},
		{"lzf", []byte{0xC3, 0}, "", rdbfile.ErrLZF},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)
			br := bufio.NewReader(bytes.NewReader(tt.input))
			actual, err := rdbfile.ReadString(br)
			if tt.err != nil {
				require.ErrorIs(err, tt.err)
			} else {
				require.Equal(tt.str, actual)
			}
		})
	}
}

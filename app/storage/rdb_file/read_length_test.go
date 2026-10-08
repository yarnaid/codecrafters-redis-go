package rdbfile_test

import (
	"bufio"
	"bytes"
	"testing"

	rdbfile "my-redis/app/storage/rdb_file"

	"github.com/stretchr/testify/require"
)

func TestReadLength(t *testing.T) {
	type res struct {
		v      uint64
		isSpec bool
		err    error
	}
	tests := []struct {
		name  string
		input []byte
		res   res
	}{
		{"simple val", []byte{0x0A}, res{v: 10, isSpec: false, err: nil}},
		{"14-bits", []byte{0x42, 0xBC}, res{v: 700, isSpec: false, err: nil}},
		{"4 bytes", []byte{0x80, 0x00, 0x00, 0x42, 0x68}, res{v: 17000, isSpec: false, err: nil}},
		{"special", []byte{0b11000000}, res{v: 0, isSpec: true, err: nil}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)
			br := bufio.NewReader(bytes.NewReader(tt.input))
			val, isSpec, err := rdbfile.ReadLength(br)
			if tt.res.err != nil {
				require.ErrorIs(err, tt.res.err)
			} else {
				require.Equal(tt.res.v, val)
				require.Equal(tt.res.isSpec, isSpec)
			}
		})
	}
}

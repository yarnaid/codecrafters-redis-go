package rdbfile_test

import (
	"bufio"
	"bytes"
	"testing"
	"time"

	rdbfile "my-redis/app/storage/rdb_file"

	"github.com/stretchr/testify/require"
)

func TestReadTimestamp(t *testing.T) {
	type conv func(br *bufio.Reader) (time.Time, error)
	tests := []struct {
		name     string
		input    []byte
		conv     conv
		expected time.Time
		err      error
	}{
		{"ms", []byte{0x15, 0x72, 0xE7, 0x07, 0x8F, 0x01, 0x00, 0x00}, rdbfile.ReadTimestampMs, time.UnixMilli(1713824559637), nil},
		{"sec", []byte{0x52, 0xED, 0x2A, 0x66}, rdbfile.ReadTimestampSec, time.Unix(1714089298, 0), nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)
			br := bufio.NewReader(bytes.NewBuffer(tt.input))
			actual, err := tt.conv(br)
			if tt.err != nil {
				require.Error(err)
			} else {
				require.Equal(tt.expected, actual)
			}
		})
	}
}

package storage_test

import (
	"fmt"
	"testing"
	"time"

	"my-redis/app/storage"

	"github.com/stretchr/testify/require"
)

func TestStreamContainer_GreaterOrGen(t *testing.T) {
	tests := []struct {
		id1, id2    string
		newStreamId storage.StreamId
		isGreater   bool
	}{
		{"1-1", "1-1", storage.StreamId{1, 1}, false},
		{"1-0", "1-1", storage.StreamId{1, 1}, false},
		{"0-1", "1-1", storage.StreamId{1, 1}, false},
		{"1-1", "0-1", storage.StreamId{0, 1}, true},
		{"1-1", "1-0", storage.StreamId{1, 0}, true},
		{"1-2", "1-1", storage.StreamId{1, 1}, true},
		{"1-2", "*", storage.StreamId{int(time.Now().UnixMilli()), 0}, false},
		{"1-2", "1-*", storage.StreamId{1, 3}, false},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s >? %s", tt.id1, tt.id2), func(t *testing.T) {
			require := require.New(t)
			id1, err := storage.ParseStreamId(tt.id1)
			require.Nil(err)
			id2, err := storage.ParseStreamId(tt.id2)
			require.Nil(err)
			actual := storage.StreamIdGreaterOrGen(&id1, &id2)
			require.Equal(tt.isGreater, actual)

			require.LessOrEqual(tt.newStreamId.Time, id2.Time)
			require.Equal(id2.Seq, tt.newStreamId.Seq)
		})
	}
}

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
		eq          bool
		newStreamId storage.StreamId
		isGreater   bool
	}{
		{"1-1", "1-1", false, storage.StreamId{1, 1}, false},
		{"1-1", "1-1", true, storage.StreamId{1, 1}, true},
		{"1-0", "1-1", false, storage.StreamId{1, 1}, false},
		{"0-1", "1-1", false, storage.StreamId{1, 1}, false},
		{"1-1", "0-1", false, storage.StreamId{0, 1}, true},
		{"1-1", "1-0", false, storage.StreamId{1, 0}, true},
		{"1-2", "1-1", false, storage.StreamId{1, 1}, true},
		{"1-2", "*", false, storage.StreamId{int(time.Now().UnixMilli()), 0}, false},
		{"1-2", "1-*", false, storage.StreamId{1, 3}, false},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s >? %s", tt.id1, tt.id2), func(t *testing.T) {
			require := require.New(t)
			id1, err := storage.ParseStreamId(tt.id1)
			require.Nil(err)
			id2, err := storage.ParseStreamId(tt.id2)
			require.Nil(err)
			actual := storage.StreamIdGreaterOrGen(&id1, &id2, tt.eq)
			require.Equal(tt.isGreater, actual)

			require.LessOrEqual(tt.newStreamId.Time, id2.Time)
			require.Equal(id2.Seq, tt.newStreamId.Seq)
		})
	}
}

func TestStreamRange(t *testing.T) {
	data := []*storage.StreamContainer{
		{Id: storage.StreamId{0, 1}},
		{Id: storage.StreamId{0, 2}},
		{Id: storage.StreamId{0, 3}},
		{Id: storage.StreamId{0, 4}},
		{Id: storage.StreamId{1, 1}},
		{Id: storage.StreamId{1, 2}},
		{Id: storage.StreamId{1, 3}},
	}

	tests := []struct {
		start, end                     string
		expectedStartId, expectedEndId storage.StreamId
		wantErr                        bool
	}{
		{"0-1", "0-3", storage.StreamId{0, 1}, storage.StreamId{0, 3}, false},
		{"0-1", "1-1", storage.StreamId{0, 1}, storage.StreamId{1, 1}, false},
		{"0-1", "0-5", storage.StreamId{0, 1}, storage.StreamId{0, 4}, false},
		{"-", "0-5", storage.StreamId{0, 1}, storage.StreamId{0, 4}, false},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("[%s:%s]", tt.start, tt.end), func(t *testing.T) {
			require := require.New(t)
			gotSlice, err := storage.StreamRange(data, tt.start, tt.end)
			require.Nil(err)
			require.Equal(tt.expectedStartId, gotSlice[0].Id)
			require.Equal(tt.expectedEndId, gotSlice[len(gotSlice)-1].Id)
		})
	}
}

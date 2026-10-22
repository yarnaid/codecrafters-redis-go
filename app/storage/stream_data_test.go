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
		newStreamId storage.StreamID
		isGreater   bool
	}{
		{"1-1", "1-1", false, storage.StreamID{1, 1}, false},
		{"1-1", "1-1", true, storage.StreamID{1, 1}, true},
		{"1-0", "1-1", false, storage.StreamID{1, 1}, false},
		{"0-1", "1-1", false, storage.StreamID{1, 1}, false},
		{"1-1", "0-1", false, storage.StreamID{0, 1}, true},
		{"1-1", "1-0", false, storage.StreamID{1, 0}, true},
		{"1-2", "1-1", false, storage.StreamID{1, 1}, true},
		{"1-2", "*", false, storage.StreamID{int(time.Now().UnixMilli()), 0}, false},
		{"1-2", "1-*", false, storage.StreamID{1, 3}, false},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s >? %s", tt.id1, tt.id2), func(t *testing.T) {
			require := require.New(t)
			id1, err := storage.ParseStreamID(tt.id1)
			require.Nil(err)
			id2, err := storage.ParseStreamID(tt.id2)
			require.Nil(err)
			actual := storage.StreamIDGreaterOrGen(&id1, &id2, tt.eq)
			require.Equal(tt.isGreater, actual)

			require.LessOrEqual(tt.newStreamId.Time, id2.Time)
			require.Equal(id2.Seq, tt.newStreamId.Seq)
		})
	}
}

func TestStreamRange(t *testing.T) {
	data := []*storage.StreamContainer{
		{ID: storage.StreamID{0, 1}},
		{ID: storage.StreamID{0, 2}},
		{ID: storage.StreamID{0, 3}},
		{ID: storage.StreamID{0, 4}},
		{ID: storage.StreamID{1, 1}},
		{ID: storage.StreamID{1, 2}},
		{ID: storage.StreamID{1, 3}},
	}

	tests := []struct {
		start, end                     string
		expectedStartId, expectedEndId storage.StreamID
		wantErr                        bool
	}{
		{"0-1", "0-3", storage.StreamID{0, 1}, storage.StreamID{0, 3}, false},
		{"0-1", "1-1", storage.StreamID{0, 1}, storage.StreamID{1, 1}, false},
		{"0-1", "0-5", storage.StreamID{0, 1}, storage.StreamID{0, 4}, false},
		{"-", "0-5", storage.StreamID{0, 1}, storage.StreamID{0, 4}, false},
		{"-", "+", storage.StreamID{0, 1}, storage.StreamID{1, 3}, false},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("[%s:%s]", tt.start, tt.end), func(t *testing.T) {
			require := require.New(t)
			gotSlice, err := storage.StreamRange(data, tt.start, tt.end)
			require.Nil(err)
			require.Equal(tt.expectedStartId, gotSlice[0].ID)
			require.Equal(tt.expectedEndId, gotSlice[len(gotSlice)-1].ID)
		})
	}
}

func TestEntriesAfter(t *testing.T) {
	data := []*storage.StreamContainer{
		{ID: storage.StreamID{0, 1}, Values: make([]storage.StreamValue, 0)},
		{ID: storage.StreamID{0, 2}, Values: make([]storage.StreamValue, 0)},
		{ID: storage.StreamID{0, 3}, Values: make([]storage.StreamValue, 0)},
	}

	tests := []struct {
		after    storage.StreamID
		expected []*storage.StreamContainer
	}{
		{storage.StreamID{0, 0}, data[:]},
		{storage.StreamID{0, 1}, data[1:]},
		{storage.StreamID{0, 2}, data[2:]},
	}

	for _, tt := range tests {
		t.Run(tt.after.String(), func(t *testing.T) {
			require := require.New(t)
			res := storage.EntriesAfter(data, tt.after)
			require.EqualValues(tt.expected, res, "got result: %v", storage.PointersToValues(res))
		})
	}
}

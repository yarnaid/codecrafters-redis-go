package storage_test

import (
	"fmt"
	"log/slog"
	"testing"
	"time"

	"my-redis/app/storage"

	"github.com/stretchr/testify/require"
)

func TestCoordinator_Xadd(t *testing.T) {
	c := storage.NewMemoryCoordinator(nil)
	key := "xadd-key"
	require := require.New(t)

	id, err := c.Xadd(key, "1-1")
	require.Nil(err)
	require.Equal("1-1", id)

	id, err = c.Xadd(key, "1-2")
	require.Nil(err)
	require.Equal("1-2", id)

	_, err = c.Xadd(key, "1-2")
	require.ErrorContains(err, "wrong streamId seq:")

	id, err = c.Xadd(key, "1-*")
	require.Nil(err, "err", err)
	require.Equal("1-3", id)

	nowMs := int(time.Now().UnixMilli())
	id, err = c.Xadd(key, "*")
	require.Nil(err)
	streamId, err := storage.ParseStreamId(id)
	require.Nil(err)
	require.GreaterOrEqual(streamId.Time, nowMs)
	require.Equal(0, streamId.Seq)

	key += "-1"
	id, err = c.Xadd(key, "0-*")
	require.Nil(err)
	require.Equal("0-1", id)

	key = "-2"
	id, err = c.Xadd(key, "*")
	require.Nil(err)
	streamId, err = storage.ParseStreamId(id)
	require.Nil(err)
	require.GreaterOrEqual(streamId.Time, nowMs)
	require.Equal(0, streamId.Seq)
}

func TestXReadStreams(t *testing.T) {
	data := []*storage.StreamContainer{
		{Id: storage.StreamId{0, 1}, Values: make([]storage.StreamValue, 0)},
		{Id: storage.StreamId{0, 2}, Values: make([]storage.StreamValue, 0)},
		{Id: storage.StreamId{0, 3}, Values: make([]storage.StreamValue, 0)},
	}
	tests := []struct {
		startId  storage.StreamId
		expected []*storage.StreamContainer
	}{
		{storage.StreamId{0, 0}, data[:]},
		{storage.StreamId{0, 1}, data[1:]},
		{storage.StreamId{0, 2}, data[2:]},
	}
	key := "xread-stream-key"
	for _, tt := range tests {
		t.Run(fmt.Sprint(tt.startId), func(t *testing.T) {
			c := storage.NewMemoryCoordinator(nil)
			require := require.New(t)

			for _, v := range data {
				_, err := c.Xadd(key, v.Id.String())
				require.Nil(err)
			}
			res, err := c.XReadStreams(key, tt.startId.String(), nil)
			require.Nil(err)
			require.EqualValues(tt.expected, res)
		})
	}
}

func TestXReadBlockStream(t *testing.T) {
	slog.Debug("Start testing `xread block streams`")
	c := storage.NewMemoryCoordinator(nil)
	key := "xread-bl-stream"
	id := storage.StreamId{1, 1}

	type result struct {
		response []*storage.StreamContainer
		err      error
	}

	startClient := func(wantWaiters int, id string) <-chan result {
		ch := make(chan result, 1) // buffered: the goroutine never leaks on send
		go func() {
			timeout := time.Duration(0) * time.Millisecond
			res, err := c.XReadStreams(key, id, &timeout)
			ch <- result{res, err}
		}()
		require.Eventually(t, func() bool { return c.WaitersStreamLen(key) >= wantWaiters },
			time.Second, time.Millisecond)
		return ch
	}

	t.Run("test 2 block clients xread stream", func(t *testing.T) {
		require := require.New(t)

		first := startClient(1, id.String())
		second := startClient(2, id.String()) // guaranteed to be queued after first

		newId := storage.StreamId{2, 2}
		_, err := c.Xadd(key, newId.String())
		require.Nil(err)

		select {
		case res := <-first:
			require.Nil(res.err)
			require.Len(res.response, 1)
			require.Equal(newId, res.response[0].Id)
		case res := <-second:
			require.Nil(res.err)
			require.Len(res.response, 1)
			require.Equal(newId, res.response[0].Id)
		case <-time.After(10 * time.Millisecond):
			require.Fail("client1 didn't receive response after. 10ms")
		}
	})

	t.Run("test $ block client xread stream", func(t *testing.T) {
		require := require.New(t)

		key = "xread-bl-$"
		res_ch := startClient(1, "$")

		newId := storage.StreamId{2, 3}
		_, err := c.Xadd(key, newId.String())
		require.Nil(err)

		select {
		case res := <-res_ch:
			require.Nil(res.err)
			require.Len(res.response, 1)
			require.Equal(newId, res.response[0].Id)
		case <-time.After(10 * time.Millisecond):
			require.Fail("client with $ didn't receive response after. 10ms")
		}
	})
}

package storage_test

import (
	"fmt"
	"log/slog"
	"testing"
	"time"

	"my-redis/app/parser"
	"my-redis/app/storage"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCoordinator_SetNGet(t *testing.T) {
	type SetVal struct {
		Key   string
		Value interface{}
		TTL   time.Duration
	}
	tests := []struct {
		name   string
		set    []SetVal
		result []bool
		err    error
	}{
		{"simple set & get", []SetVal{{"123", 123, 0}}, []bool{true}, nil},
		{"set & get 2 vals", []SetVal{{"123", 123, 0}, {"124", 124, 0}}, []bool{true, true}, nil},
		{"set & get twice", []SetVal{{"123", 123, 0}, {"123", 123, 0}}, []bool{true, true}, nil},
		{"set & get with TTL", []SetVal{{"123", 123, 100 * time.Millisecond}}, []bool{true}, nil},
	}

	assert := assert.New(t)
	require := require.New(t)
	slog.SetLogLoggerLevel(slog.LevelDebug)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slog.Debug(tt.name)
			s := storage.NewMemoryCoordinator(nil)
			for i, v := range tt.set {
				ret := s.Set(v.Key, v.Value, v.TTL)
				assert.Equal(tt.result[i], ret)

				val, ok := s.Get(v.Key)
				assert.True(ok)
				assert.Equal(v.Value, val)
				if v.TTL > 0 {
					time.Sleep(v.TTL + time.Millisecond)
					val, ok = s.Get(v.Key)
					require.False(ok)
					require.Nil(val)
				}

				// assert.True(reflect.DeepEqual(v.Value, val), "%v != %v", v.Value, val)
			}
		})
	}

	_ = tests
}

func TestCoordinator_Append(t *testing.T) {
	tests := []struct {
		name   string
		values []parser.Serializable
	}{
		{"simple", []parser.Serializable{}},
	}
	key := "key"

	require := require.New(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := storage.NewMemoryCoordinator(nil)
			for _, v := range tt.values {
				res := s.Append(key, v)
				require.Equal(1, res)
			}
			res, _ := s.Get(key)
			res_arr, _ := res.([]interface{})
			require.Equal(len(tt.values), len(res_arr))
		})
	}
}

func TestCoordinator_LRange(t *testing.T) {
	l := 7
	data := make(parser.Array[parser.Serializable], l)
	for i := range l {
		data[i] = parser.Int(i)
	}

	key := "key"
	tests := []struct {
		start, end int
		res        parser.Array[parser.Serializable]
	}{
		{0, l - 1, data[:l]},
		{0, -1, data[:l]},
		{0, -2, data[:l-1]},
		{0, 0, data[:1]},
		{0, 1, data[:2]},
		{1, -2, data[1 : l-1]},
		{-4, -2, data[l-4 : l-1]},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d, %d", tt.start, tt.end), func(t *testing.T) {
			slog.Debug("starting test", "start", tt.start, "end", tt.end)
			require := require.New(t)
			c := storage.NewMemoryCoordinator(nil)
			c.Append(key, data...)
			// store, ok := c.Get(key)
			// require.True(ok, "array in storage not found!")
			// slog.Debug("starting execution", "data", data, "storage", store)

			actual := c.LRange(key, tt.start, tt.end)
			require.EqualValues(tt.res, actual)
		})
	}
}

func TestCoordinator_RPush(t *testing.T) {
	tests := []struct {
		name     string
		values   []parser.Serializable
		expected []parser.Serializable
	}{
		{"1 int", []parser.Serializable{parser.Int(1)}, []parser.Serializable{parser.Int(1)}},
		{"2 int", []parser.Serializable{parser.Int(1), parser.Int(2)}, []parser.Serializable{parser.Int(1), parser.Int(2)}},
	}
	key := "key"

	slog.SetLogLoggerLevel(slog.LevelDebug)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slog.Debug("starting rpush test")
			require := require.New(t)
			c := storage.NewMemoryCoordinator(nil)
			res := c.Append(key, tt.values...)
			require.Equal(len(tt.values), res)

			arr, ok := c.Get(key)
			require.True(ok)
			require.EqualValues(tt.expected, arr)
		})
	}
}

func TestCoordinator_LPop2Clients(t *testing.T) {
	c := storage.NewMemoryCoordinator(nil)
	key := "coord-bl-pop-key"
	value := "coord-bl-pop-value"

	type result struct {
		res parser.Serializable
		err error
	}
	start_client := func(r chan result) {
		res, err := c.BLPop(key, 0)
		r <- result{res, err}
	}
	result_ch1 := make(chan result)
	result_ch2 := make(chan result)
	go start_client(result_ch1)
	go start_client(result_ch2)

	t.Run("2 clients test", func(t *testing.T) {
		require := require.New(t)
		require.Eventually(func() bool { return c.WaitersLen(key) > 0 }, time.Second, 50*time.Millisecond)

		n := c.Append(key, parser.BulkString(value))
		require.Equal(0, n)

		select {
		case res := <-result_ch1:
			require.Nil(res.err)
			require.Equal(parser.BulkString(value), res.res)
		case res := <-result_ch2:
			require.Fail("second client got unexpected response", "res", res)
		case <-time.After(10 * time.Millisecond):
			require.Fail("client1 didn't receive response after. 10ms")
		}
	})
}

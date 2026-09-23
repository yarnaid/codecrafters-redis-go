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
		Kind  storage.ValueKind
	}
	tests := []struct {
		name   string
		set    []SetVal
		result []bool
		err    error
	}{
		{"simple set & get", []SetVal{{"123", 123, 0, storage.KindString}}, []bool{true}, nil},
		{"set & get 2 vals", []SetVal{{"123", 123, 0, storage.KindString}, {"124", 124, 0, storage.KindString}}, []bool{true, true}, nil},
		{"set & get twice", []SetVal{{"123", 123, 0, storage.KindString}, {"123", 123, 0, storage.KindString}}, []bool{true, true}, nil},
		{"set & get with TTL", []SetVal{{"123", 123, 100 * time.Millisecond, storage.KindString}}, []bool{true}, nil},
	}

	assert := assert.New(t)
	require := require.New(t)
	slog.SetLogLoggerLevel(slog.LevelDebug)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slog.Debug(tt.name)
			s := storage.NewMemoryCoordinator(nil)
			for i, v := range tt.set {
				ret := s.Set(v.Key, v.Value, v.TTL, v.Kind)
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
		// {"0", []parser.Serializable{}},
		{"int 1", []parser.Serializable{parser.Int(0)}},
		{"int 2", []parser.Serializable{parser.Int(0), parser.Int(1)}},
		{"int 3", []parser.Serializable{parser.Int(0), parser.Int(1), parser.Int(2)}},
	}
	key := "coord-append-key"

	require := require.New(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := storage.NewMemoryCoordinator(nil)
			s.Append(key, tt.values...)
			res, ok := s.Get(key)
			require.True(ok, "cannot get key %v", key)
			res_arr, ok := res.([]parser.Serializable)
			require.True(ok, "cannot convert value %v of %T", res, res)
			vv, _ := s.Backend().Get(key)
			require.Equal(len(tt.values), len(res_arr), "incorrect array len %v\nval %v", res_arr, vv)
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

func TestCoordinator_BLPop2Clients(t *testing.T) {
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
		require.Equal(1, n)

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

func TestCoordinator_LLen(t *testing.T) {
	key := "coord-llen-key"
	tests := []struct {
		values   []parser.Serializable
		expected int
	}{
		{[]parser.Serializable{parser.Int(0)}, 1},
		{[]parser.Serializable{parser.Int(0), parser.Int(1)}, 2},
		{[]parser.Serializable{parser.Int(0), parser.Int(1), parser.Int(2)}, 3},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprint(tt.expected), func(t *testing.T) {
			require := require.New(t)
			c := storage.NewMemoryCoordinator(nil)
			c.Append(key, tt.values...)
			stored_inter, ok := c.Get(key)
			require.True(ok)
			stored, ok := stored_inter.([]parser.Serializable)
			require.True(ok, "got value %v of %T", stored_inter, stored_inter)
			require.Equal(tt.expected, len(stored), "got %v (%d) of type %t", stored, len(stored), stored)
			res := c.LLen(key)
			require.Equal(tt.expected, res, "incorrect array=%v", stored)
		})
	}
}

func TestCoordinator_Type(t *testing.T) {
	c := storage.NewMemoryCoordinator(nil)
	require := require.New(t)
	c.Set("str", "str", 0, storage.KindString)
	v, err := c.Type("str")
	require.Nil(err)
	require.Equal(storage.KindString, v)

	c.Append("list", []parser.Serializable{parser.BulkString("str")}...)
	v, err = c.Type("list")
	require.Nil(err)
	require.Equal(storage.KindList, v)
}

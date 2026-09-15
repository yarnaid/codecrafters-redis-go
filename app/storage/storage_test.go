package storage_test

import (
	"log/slog"
	"testing"
	"time"

	"my-redis/app/storage"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStorage_SetNGet(t *testing.T) {
	type SetVal struct {
		Key   string
		Value interface{}
		TTL   int
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
		{"set & get with TTL", []SetVal{{"123", 123, 100}}, []bool{true}, nil},
	}

	assert := assert.New(t)
	// require := require.New(t)
	slog.SetLogLoggerLevel(slog.LevelDebug)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slog.Debug(tt.name)
			s := storage.NewStorage()
			for i, v := range tt.set {
				ret := s.Set(v.Key, v.Value, v.TTL)
				assert.Equal(tt.result[i], ret)

				val, ok := s.Get(v.Key)
				assert.True(ok)
				assert.Equal(v.Value, val)
				if v.TTL > 0 {
					time.Sleep(time.Duration(v.TTL+10) * time.Microsecond)
					val, ok = s.Get(v.Key)
					assert.False(ok)
					assert.Nil(val)
				}

				// assert.True(reflect.DeepEqual(v.Value, val), "%v != %v", v.Value, val)
			}
		})
	}

	_ = tests
}

func TestStorage_Append(t *testing.T) {
	tests := []struct {
		name   string
		values []interface{}
	}{
		{"simple", []interface{}{}},
	}
	key := "key"

	require := require.New(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := storage.NewStorage()
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

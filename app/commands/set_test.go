package commands_test

import (
	"testing"
	"time"

	. "my-redis/app/commands"
	"my-redis/app/parser"
	"my-redis/app/storage"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetCommand(t *testing.T) {
	type set_args struct {
		Key   string
		Value interface{}
		TTL   time.Duration

		ResVal parser.Serializable
		Error  error
	}
	tests := []struct {
		name string
		args []set_args
	}{
		{"simple", []set_args{{"123", 123, 0, parser.SimpleString("OK"), nil}}},
		{"simple twice", []set_args{{"123", 123, 0, parser.SimpleString("OK"), nil}, {"123", 123, 0, parser.SimpleString("OK"), nil}}},
		{"simple with ttl", []set_args{{"123", 123, 100 * time.Millisecond, parser.SimpleString("OK"), nil}}},
	}

	assert := assert.New(t)
	require := require.New(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := storage.NewMemoryCoordinator(nil)
			expVersion := len(tt.args)
			var wantTimeout bool
			for _, ttt := range tt.args {
				cmd := SetCommand{s, ttt.Key, ttt.Value, ttt.TTL}
				res, err := cmd.Execute()
				assert.Equal(ttt.ResVal, res)
				assert.Nil(err)
				s_val, ok := s.Get(ttt.Key)
				assert.True(ok)
				assert.Equal(ttt.Value, s_val)

				if ttt.TTL > 0 {
					time.Sleep(ttt.TTL + time.Microsecond)
					s_val, ok = s.Get(ttt.Key)
					require.False(ok, "value must be missing after ttl")
					require.Equal(nil, s_val)
					wantTimeout = true
				}
			}
			version := s.GetVersions(tt.args[0].Key)[0]
			if !wantTimeout {
				require.Equal(expVersion, version)
			}
		})
	}
}

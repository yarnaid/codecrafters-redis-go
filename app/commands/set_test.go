package commands_test

import (
	"testing"
	"time"

	. "my-redis/app/commands"
	"my-redis/app/parser"
	"my-redis/app/storage"

	"github.com/stretchr/testify/assert"
)

func TestSetCommand(t *testing.T) {
	type set_args struct {
		Key    string
		Value  interface{}
		TTL_MS int

		ResVal parser.Serializable
		Error  error
	}
	tests := []struct {
		name string
		args []set_args
	}{
		{"simple", []set_args{{"123", 123, 0, parser.SimpleString("OK"), nil}}},
		{"simple twice", []set_args{{"123", 123, 0, parser.SimpleString("OK"), nil}, {"123", 123, 0, parser.SimpleString("OK"), nil}}},
		{"simple with ttl", []set_args{{"123", 123, 100, parser.SimpleString("OK"), nil}}},
	}

	assert := assert.New(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := storage.NewStorage()
			for _, ttt := range tt.args {
				cmd := SetCommand{s, ttt.Key, ttt.Value, ttt.TTL_MS}
				res, err := cmd.Execute()
				assert.Equal(ttt.ResVal, res)
				assert.Nil(err)
				s_val, ok := s.Get(ttt.Key)
				assert.True(ok)
				assert.Equal(ttt.Value, s_val)

				if ttt.TTL_MS > 0 {
					time.Sleep(time.Duration(ttt.TTL_MS+10) * time.Microsecond)
					s_val, ok = s.Get(ttt.Key)
					assert.False(ok)
					assert.Equal(nil, s_val)
				}
			}
		})
	}
}

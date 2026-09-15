package commands_test

import (
	"log/slog"
	"testing"
	"time"

	"my-redis/app/commands"
	"my-redis/app/parser"
	"my-redis/app/storage"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCommand_Execute(t *testing.T) {
	type kv struct {
		Key     string
		Value   interface{}
		TTL_MS  int
		wantErr bool
	}
	tests := []struct {
		name  string // description of this test case
		store []kv
	}{
		{"simple", []kv{{"123", 123, 0, false}}},
		{"ttl", []kv{{"123", 123, 100, false}}},
		{"twice", []kv{{"123", 123, 0, false}, {"123", 123, 0, false}}},
		{"twice diff", []kv{{"123", 123, 0, false}, {"str", "str", 0, false}}},
	}
	assert := assert.New(t)
	require := require.New(t)
	slog.SetLogLoggerLevel(slog.LevelDebug)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := storage.NewStorage()
			for _, v := range tt.store {
				slog.Debug("Want set", "key", v.Key, "value", v.Value)
				require.True(s.Set(v.Key, v.Value, v.TTL_MS))
				slog.Debug("Success set", "key", v.Key, "value", v.Value)

				p := commands.GetCommand{s, v.Key}
				got, gotErr := p.Execute()
				slog.Debug("GET cmd executed", "got", got, "gotErr", gotErr)
				if !v.wantErr {
					assert.Nil(gotErr)
					switch vv := v.Value.(type) {
					case int:
						assert.Equal(parser.Int(vv), got)
					case string:
						assert.Equal(parser.BulkString(vv), got)
					}

					if v.TTL_MS > 0 {
						time.Sleep(time.Duration(v.TTL_MS+10) * time.Microsecond)
						p := commands.GetCommand{s, v.Key}
						got, gotErr = p.Execute()
						require.Nil(gotErr)
						require.Equal(parser.BulkNullString(""), got)

					}

				} else {
					assert.NotNil(gotErr)
					assert.Nil(got)
				}
			}
		})
	}
}

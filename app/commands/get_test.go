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
		Value   parser.BulkString
		TTL     time.Duration
		Kind    storage.ValueKind
		wantErr bool
	}
	tests := []struct {
		name  string // description of this test case
		store []kv
	}{
		{"simple", []kv{{"123", "123", 0, storage.KindString, false}}},
		{"ttl", []kv{{"124", "124", 100 * time.Millisecond, storage.KindString, false}}},
		{"twice", []kv{{"125", "125", 0, storage.KindString, false}, {"126", "126", 0, storage.KindString, false}}},
		{"twice diff", []kv{{"127", "127", 0, storage.KindString, false}, {"str", "str", 0, storage.KindString, false}}},
	}
	assert := assert.New(t)
	require := require.New(t)
	slog.SetLogLoggerLevel(slog.LevelDebug)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := storage.NewMemoryCoordinator(nil)
			for _, v := range tt.store {
				// slog.Debug("Want set", "key", v.Key, "value", v.Value)
				require.True(s.Set(v.Key, v.Value, v.TTL, v.Kind))
				// slog.Debug("Success set", "key", v.Key, "value", v.Value)

				p := commands.GetCommand{s, v.Key}
				got, gotErr := p.Execute()
				// slog.Debug("GET cmd executed", "got", got, "gotErr", gotErr)
				if !v.wantErr {
					require.NoError(gotErr)
					// vv, ok := v.Value.(parser.BulkString)
					// require.True(ok, "incorrect type of value %T", v.Value)
					require.Equal(parser.BulkString(v.Value), got)

					if v.TTL > 0 {
						time.Sleep(v.TTL + time.Millisecond)
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

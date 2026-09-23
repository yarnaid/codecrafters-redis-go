package commands_test

import (
	"log/slog"
	"testing"
	"time"

	. "my-redis/app/commands"
	"my-redis/app/parser"
	"my-redis/app/storage"

	"github.com/stretchr/testify/require"
)

type (
	Be storage.Backend
)

func TestBLPopSimple(t *testing.T) {
	key := "key"
	tests := []struct {
		name        string
		backend     func() Be
		timeout     time.Duration
		result      string
		wantTimeout bool
		wantErr     bool
	}{
		{"simple without timeout", func() Be {
			data := make(map[string]*storage.Value)
			data[key] = &storage.Value{Value: []parser.Serializable{parser.BulkString("123")}}
			return storage.NewMemoryBackend(data)
		}, time.Millisecond * 50, "123", false, false},
		{"simple timeout", func() Be { return &storage.MemoryBackend{} }, time.Millisecond * 50, "", true, false},
	}

	type result struct {
		res parser.Serializable
		err error
	}

	slog.SetLogLoggerLevel(slog.LevelDebug)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)
			be := tt.backend().(*storage.MemoryBackend)
			coordinator := storage.NewMemoryCoordinator(be)

			res_ch := make(chan result, 1)
			go func() {
				cmd_res, err := BLPop{coordinator, "key", tt.timeout}.Execute()
				res_ch <- result{cmd_res, err}
				slog.Debug("[test][blpop] value received", "res", cmd_res, "err", err)
			}()

			select {
			case res := <-res_ch:
				slog.Debug("[test][blpop] value received", "res", res)
				if !tt.wantErr {
					require.Nil(res.err)
					if !tt.wantTimeout {
						require.EqualValues(parser.Array[parser.Serializable]{parser.BulkString("key"), parser.BulkString(tt.result)}, res.res)
					} else {
						require.EqualValues(parser.NullArray(0), res.res)
					}
				} else {
					require.NotNil(res.err)
				}
			case <-time.After(tt.timeout + 50*time.Millisecond):
				if tt.wantTimeout {
					t.Fatal("Timeout didn't fire")
				}
			}
		})
	}
}

func TestBLPop2Clients(t *testing.T) {
	c := storage.NewMemoryCoordinator(nil)

	key := "blpop-key"
	value := "blpop-hello"
	type result struct {
		res parser.Serializable
		err error
	}
	result_ch1 := make(chan result, 1)
	result_ch2 := make(chan result, 1)
	client_f := func(ch chan result) {
		cmd := BLPop{S: c, Key: key, Timeout: 0}
		slog.Debug("[testing][blpop-command] command created")
		res, err := cmd.Execute()
		slog.Debug("[testing][blpop-command] execution finished", "res", res, "err", err)
		ch <- result{res, err}
	}
	go client_f(result_ch1)
	time.Sleep(10 * time.Microsecond) // dirty hack
	go client_f(result_ch2)

	time.Sleep(10 * time.Millisecond) // dirty hack

	t.Run("test2clients", func(t *testing.T) {
		require := require.New(t)
		push := &RPushCommand{c, key, []parser.Serializable{parser.BulkString(value)}}
		res, err := push.Execute()
		require.Nil(err)
		require.Equal(parser.Int(0), res)

		select {
		case <-time.After(10 * time.Millisecond):
			require.Fail("[testing][blpop-command] No response got in 10ms!")
		case <-result_ch2:
			require.Fail("[testing][blpop-command] Second client got response")
		case cl1_res := <-result_ch1:
			require.Nil(cl1_res.err)
			require.Equal(parser.Array[parser.Serializable]{parser.BulkString(key), parser.BulkString(value)}, cl1_res.res)
		}
		slog.Debug("[testing][blpop-command] BLPOP 2 Clients test finished")
	})
}

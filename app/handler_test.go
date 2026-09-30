package main_test

import (
	"net"
	"testing"

	. "my-redis/app"
	"my-redis/app/parser"
	"my-redis/app/storage"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler(t *testing.T) {
	tests := []struct {
		name   string
		input  []byte
		output []byte
	}{
		{"type", []byte("*2\r\n$4\r\nTYPE\r\n$4\r\nNONE\r\n"), []byte("+none\r\n")},
		{"ping", []byte("*1\r\n$4\r\nPING\r\n"), []byte("+PONG\r\n")},
		{"echo", []byte("*2\r\n$4\r\nECHO\r\n$11\r\nhello world\r\n"), []byte("$11\r\nhello world\r\n")},
		{"echo", []byte("*2\r\n$4\r\nECHO\r\n$9\r\npineapple\r\n"), []byte("$9\r\npineapple\r\n")},
		{"set", []byte("*3\r\n$3\r\nSET\r\n$3\r\nKEY\r\n$3\r\n123\r\n"), []byte("+OK\r\n")},
		{"rpush", []byte("*3\r\n$5\r\nRPUSH\r\n$3\r\nKEY\r\n$3\r\n123\r\n"), []byte(":1\r\n")},
		{"set with PX", []byte("*5\r\n$3\r\nSET\r\n$10\r\nstrawberry\r\n$6\r\nbanana\r\n$2\r\nPX\r\n$3\r\n100\r\n"), []byte("+OK\r\n")},
		{"multi", []byte("*1\r\n$5\r\nMULTI\r\n"), []byte("+OK\r\n")},
		{"exec", []byte("*1\r\n$4\r\nEXEC\r\n"), []byte("-ERR EXEC without MULTI\r\n")},
		// {"get", []byte("*2\r\n$3\r\nGET\r\n$3\r\nKEY\r\n"), []byte("-\r\npineapple\r\n")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)
			client := initClient()

			_, err := client.Write(tt.input)
			require.NoError(err)

			buf := make([]byte, 512)
			n, err := client.Read(buf)
			require.NoError(err)

			response := buf[:n]
			assert.Equal(t, string(tt.output), string(response))
		})
	}
}

func initClient() net.Conn {
	client, server := net.Pipe()
	handler := Handler{
		Conn:        server,
		Coordinator: storage.NewMemoryCoordinator(nil),
	}
	go handler.HandleConnection()
	return client
}

func stringsToCmd(input []string) []byte {
	a := make(parser.Array[parser.BulkString], len(input))
	for i := range len(input) {
		a[i] = parser.BulkString(input[i])
	}
	res, _ := a.Serialize()
	return res
}

func TestHandlerSeq(t *testing.T) {
	type cmdAndRes struct {
		cmd []string
		res func() []byte
	}
	tests := []struct {
		name     string
		commands []cmdAndRes
	}{
		{"echo", []cmdAndRes{{cmd: []string{"ECHO", "123"}, res: func() []byte { return []byte("$3\r\n123\r\n") }}}},
		{"empty tx", []cmdAndRes{
			{cmd: []string{"MULTI"}, res: func() []byte { return []byte("+OK\r\n") }},
			{cmd: []string{"SET", "key", "value"}, res: func() []byte { return []byte("+QUEUED\r\n") }},
			{cmd: []string{"EXEC"}, res: func() []byte { res, _ := parser.Array[parser.Serializable]{}.Serialize(); return res }},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := initClient()
			require := require.New(t)
			for _, cmd := range tt.commands {
				_, err := client.Write(stringsToCmd(cmd.cmd))
				require.NoError(err)

				buf := make([]byte, 512)
				n, err := client.Read(buf)
				require.NoError(err)

				response := buf[:n]
				require.NoError(err)
				require.Equal(string(cmd.res()), string(response))

			}
		})
	}
}

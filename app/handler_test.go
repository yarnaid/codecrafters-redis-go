package main_test

import (
	"net"
	"testing"

	. "my-redis/app"
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
		{"ping", []byte("*1\r\n$4\r\nPING\r\n"), []byte("+PONG\r\n")},
		{"echo", []byte("*2\r\n$4\r\nECHO\r\n$11\r\nhello world\r\n"), []byte("$11\r\nhello world\r\n")},
		{"echo", []byte("*2\r\n$4\r\nECHO\r\n$9\r\npineapple\r\n"), []byte("$9\r\npineapple\r\n")},
		{"set", []byte("*3\r\n$3\r\nSET\r\n$3\r\nKEY\r\n$3\r\n123\r\n"), []byte("+OK\r\n")},
		{"set with PX", []byte("*5\r\n$3\r\nSET\r\n$10\r\nstrawberry\r\n$6\r\nbanana\r\n$2\r\nPX\r\n$3\r\n100\r\n"), []byte("+OK\r\n")},
		// {"get", []byte("*2\r\n$3\r\nGET\r\n$3\r\nKEY\r\n"), []byte("-\r\npineapple\r\n")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, server := net.Pipe()
			handler := Handler{
				Conn: server,
				Strg: storage.NewStorage(),
			}
			go handler.HandleConnection()

			_, err := client.Write(tt.input)
			require.NoError(t, err)

			buf := make([]byte, 512)
			n, err := client.Read(buf)
			require.NoError(t, err)

			response := buf[:n]
			assert.Equal(t, tt.output, response)
		})
	}
}

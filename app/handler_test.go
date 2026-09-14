package main

import (
	"testing"

	"net"

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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, server := net.Pipe()
			handler := Handler{}
			handler.SetConn(server)
			go handler.handle_connection()

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

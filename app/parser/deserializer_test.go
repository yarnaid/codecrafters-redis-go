package parser

import (
	"log/slog"
	"testing"

	"reflect"

	"github.com/stretchr/testify/assert"
)

func TestDeserialization(t *testing.T) {
	// slog.SetLogLoggerLevel(slog.LevelDebug)
	tests := []struct {
		name     string
		input    []byte
		expected interface{}
		error    bool
	}{
		{"int", []byte(":42\r\n"), Int(42), false},
		{"simple string", []byte("+Hello\r\n"), SimpleString("Hello"), false},
		{"bulk string", []byte("$5\r\nHello\r\n"), BulkString("Hello"), false},
		{"simple error", []byte("-Error\r\n"), SimpleError("Error"), false},
		{"invalid input", []byte("Invalid\r\n"), nil, true},
		{"array of integers", []byte("*3\r\n:1\r\n:2\r\n:3\r\n"), []Serializable{Int(1), Int(2), Int(3)}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Deserialize(tt.input)
			slog.Debug("Deserialization test", "name", tt.name, "result", result, "error", err)
			if (err != nil) != tt.error {
				assert.Fail(t, "Deserialize() error = %v, wantErr %v", err, tt.error)
				return
			}
			if !tt.error && !reflect.DeepEqual(result, tt.expected) {
				assert.Fail(t, "Deserialize() = %v, want %v", result, tt.expected)
			}
		})
	}
}

package parser_test

import (
	"testing"

	. "my-redis/app/parser"
)

func TestSerialization(t *testing.T) {
	tests := []struct {
		name     string
		input    Serializable
		expected []byte
		error    bool
	}{
		{"int", Int(42), []byte(":42\r\n"), false},
		{"simple string", SimpleString("Hello"), []byte("+Hello\r\n"), false},
		{"bulk string", BulkString("Hello"), []byte("$5\r\nHello\r\n"), false},
		{"simple error", SimpleError("Error"), []byte("-Error\r\n"), false},
		{"simple string with CRLF", SimpleString("Hello\r\nWorld"), nil, true},
		{"array of simple strings", Array[Serializable]{SimpleString("Hello"), SimpleString("World")}, []byte("*2\r\n+Hello\r\n+World\r\n"), false},
		{"array of mixed types", Array[Serializable]{Int(42), SimpleString("Hello"), BulkString("World")}, []byte("*3\r\n:42\r\n+Hello\r\n$5\r\nWorld\r\n"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Serialize(tt.input)
			if (err != nil) != tt.error {
				t.Errorf("Serialize() error = %v, wantErr %v", err, tt.error)
				return
			}
			if !tt.error && string(result) != string(tt.expected) {
				t.Errorf("Serialize() = %v, want %v", result, tt.expected)
			}
		})
	}
}

package commands_test

import (
	"testing"

	"my-redis/app/commands"
	"my-redis/app/parser"
	"my-redis/app/storage"

	"github.com/stretchr/testify/require"
)

func TestRPushCommand(t *testing.T) {
	tests := []struct {
		name   string
		values []interface{}
	}{
		{"single value", []interface{}{1}},
		{"multiple values", []interface{}{1, 2, 3}},
	}
	// assert := assert.New(t)
	require := require.New(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := storage.NewStorage()
			cmd := commands.RPushCommand{s, "test", tt.values}
			got, gotErr := cmd.Execute()
			require.Nil(gotErr)
			require.Equal(parser.Int(len(tt.values)), got)
		})
	}
}

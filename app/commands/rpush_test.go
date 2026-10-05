package commands_test

import (
	"context"
	"testing"

	"my-redis/app/commands"
	. "my-redis/app/parser"
	"my-redis/app/storage"

	"github.com/stretchr/testify/require"
)

func TestRPushCommand(t *testing.T) {
	tests := []struct {
		name   string
		values []Serializable
	}{
		{"single value", []Serializable{Int(1)}},
		{"multiple values", []Serializable{Int(1), Int(2), Int(3)}},
	}
	// assert := assert.New(t)
	require := require.New(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := storage.NewMemoryCoordinator(nil)
			env := &commands.Env{Coordinator: s}
			cmd := commands.RPushCommand{"test", tt.values}
			got, gotErr := cmd.Execute(context.Background(), env, nil)
			require.Nil(gotErr)
			require.Equal(Int(len(tt.values)), got)
		})
	}
}

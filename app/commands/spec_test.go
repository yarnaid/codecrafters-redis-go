package commands_test

import (
	"testing"

	"my-redis/app/commands"

	"github.com/stretchr/testify/require"
)

func TestSpec_Validate(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		spec    commands.Spec
		wantErr bool
	}{
		{"valid", []string{"SET", "key", "value"}, commands.Spec{Name: "SET", MinArgs: 2, MaxArgs: 2}, false},
		{"too few args", []string{"SET", "key"}, commands.Spec{Name: "SET", MinArgs: 2, MaxArgs: 2}, true},
		{"too many args", []string{"SET", "key", "value", "extra"}, commands.Spec{Name: "SET", MinArgs: 2, MaxArgs: 2}, true},
		{"no max args", []string{"SET", "key", "value"}, commands.Spec{Name: "SET", MinArgs: 2, MaxArgs: -1}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)
			gotErr := tt.spec.Validate(tt.args[1:])
			if tt.wantErr {
				require.Error(gotErr)
			} else {
				require.NoError(gotErr)
			}
		})
	}
}

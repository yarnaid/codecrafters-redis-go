package commands_test

import (
	"fmt"
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

func TestFlag_Propagate(t *testing.T) {
	tests := []struct {
		f   commands.Flags
		res bool
	}{
		{0, true},
		{commands.FlagNotPropagate, false},
		{commands.FlagBlocking, true},
		{commands.FlagNotPropagate | commands.FlagBlocking, false},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprint(tt.f), func(t *testing.T) {
			require := require.New(t)
			require.Equal(tt.res, tt.f.Propagate())
		})
	}
}

func TestFlags_Has(t *testing.T) {
	tests := []struct {
		name    string
		flag    commands.Flags
		another commands.Flags
		has     bool
	}{
		{"one has one", 1, 1, true},
		{"2 has one", 3, 1, true},
		{"2 has one", 3, 2, true},
		{"1 doesn't have one", 4, 2, false},
		{"allowed in subs", commands.FlagNotPropagate & commands.FlagAllowedInSubs, commands.FlagAllowedInSubs, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)
			actual := tt.flag.Has(tt.another)
			require.Equal(tt.has, actual)
		})
	}
}

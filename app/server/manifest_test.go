package server_test

import (
	"testing"

	"my-redis/app/server"

	"github.com/stretchr/testify/require"
)

func TestGetNextAOFIncr(t *testing.T) {
	tests := []struct {
		name     string
		names    []string
		base     string
		expected int
	}{
		{"simple", []string{"A.aof.1.incr.aof"}, "A.aof", 2},
		{"2 old", []string{"A.aof.1.incr.aof", "A.aof.2.incr.aof"}, "A.aof", 3},
		{"the new", []string{}, "A.aof", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)
			act, err := server.GetNextAOFIncr(tt.names, tt.base)
			require.NoError(err)
			require.Equal(tt.expected, act)
		})
	}
}

func TestNewAOFName(t *testing.T) {
	require := require.New(t)
	res, err := server.NewAOFFileName("/", "a.aof")
	require.NoError(err)
	require.Equal("a.1.incr.aof", res)
}

func TestParseManifest(t *testing.T) {
	tests := []struct {
		name string
		s    string
		m    server.Manifest
		err  error
	}{
		{"simple", "file test.aof.1.incr.aof seq 1 type i", server.Manifest{"test.aof.1.incr.aof", 1, "i"}, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)
			actual, err := server.ParseManifest(tt.s)
			if tt.err != nil {
				require.Error(err)
			} else {
				require.NoError(err)
				require.Equal(tt.m, actual)
			}
		})
	}
}

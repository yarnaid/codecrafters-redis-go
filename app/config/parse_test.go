package config_test

import (
	"testing"

	"my-redis/app/config"

	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	require := require.New(t)
	cfg, err := config.ParseConfig([]string{})
	require.NoError(err)
	require.Equal("no", cfg.AppendOnly)
	m, err := cfg.ToMap()
	require.NoError(err)
	v, ok := m["appendonly"]
	require.True(ok, "%v\n%v", cfg, m)
	require.Equal("no", v)
}

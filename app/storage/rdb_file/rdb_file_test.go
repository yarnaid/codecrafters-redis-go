package rdbfile_test

import (
	"bytes"
	_ "embed"
	"testing"

	rdbfile "my-redis/app/storage/rdb_file"

	"github.com/stretchr/testify/require"
)

//go:embed testdata/dump_test.rdb
var sampleRDB []byte

func TestParse(t *testing.T) {
	require := require.New(t)
	require.Greater(len(sampleRDB), 9)
	snap, err := rdbfile.Parse(bytes.NewReader(sampleRDB))
	require.NoError(err)
	require.NotNil(snap)
	require.Equal(15, snap.Version)
	require.Equal("8.10.2", snap.Metadata["redis-ver"])
	require.Equal("bar", snap.Data["foo"].Value)
	require.True(snap.Data["foo"].ExpireAt.IsZero())
	require.Len(snap.Data, 1)
}

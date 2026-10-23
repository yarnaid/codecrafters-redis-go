package user

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_flagSet_List(t *testing.T) {
	tests := []struct {
		name  string // description of this test case
		flags []flag
		want  []string
	}{
		{"empty", []flag{}, []string{}},
		{"nopass", []flag{FlagNoPass}, []string{"nopass"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)
			var f flagSet
			for _, flg := range tt.flags {
				f |= flg.fs()
			}
			got := f.List()
			require.Equal(tt.want, got)
		})
	}
}

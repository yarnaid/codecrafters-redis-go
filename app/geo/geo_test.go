package geo_test

import (
	"testing"

	"my-redis/app/geo"

	"github.com/stretchr/testify/require"
)

type testCase struct {
	Name      string
	Lat, Long float64
	Score     uint64
}

var testCases = []testCase{
	{"Bangkok", 13.7220, 100.5252, 3962257306574459},
	{"Beijing", 39.9075, 116.3972, 4069885364908765},
	{"Berlin", 52.5244, 13.4105, 3673983964876493},
	{"Copenhagen", 55.6759, 12.5655, 3685973395504349},
	{"New Delhi", 28.6667, 77.2167, 3631527070936756},
	{"Kathmandu", 27.7017, 85.3206, 3639507404773204},
	{"London", 51.5074, -0.1278, 2163557714755072},
	{"New York", 40.7128, -74.0060, 1791873974549446},
	{"Paris", 48.8534, 2.3488, 3663832752681684},
	{"Sydney", -33.8688, 151.2093, 3252046221964352},
	{"Tokyo", 35.6895, 139.6917, 4171231230197045},
	{"Vienna", 48.2064, 16.3707, 3673109836391743},
}

func TestLoc_Score(t *testing.T) {
	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {
			require := require.New(t)
			l := geo.NewLoc(tt.Lat, tt.Long, 0)
			got := l.Score()
			require.Equal(tt.Score, got)
		})
	}
}

func TestLoc_Coord(t *testing.T) {
	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {
			require := require.New(t)
			l := geo.NewLoc(0, 0, tt.Score)
			gotLat, got2Long := l.Coord()

			require.InDelta(tt.Lat, gotLat, 1e-5)
			require.InDelta(tt.Long, got2Long, 1e-5)
		})
	}
}

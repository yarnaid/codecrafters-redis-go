package geo

import (
	"fmt"
	"math"
)

const (
	MinLong   = -180.
	MaxLong   = 180.
	MinLat    = -85.05112878
	MaxLat    = 85.05112878
	LatRange  = MaxLat - MinLat
	LongRange = MaxLong - MinLong
	rEarth    = 6372797.560856
)

type (
	grad float64
	rad  float64
)

func (g grad) toRad() rad {
	return rad(float64(g) / 180. * math.Pi)
}

func (g grad) v() float64 {
	return float64(g)
}

func (r rad) v() float64 {
	return float64(r)
}

type Loc struct {
	lat, long grad
	score     uint64
}

func NewLoc(lat, long float64, score uint64) *Loc {
	res := &Loc{lat: grad(lat), long: grad(long), score: score}
	if score != 0 {
		res.lat, res.long = res.Coord()
	}
	return res
}

func (l *Loc) String() string {
	l1, l2 := l.Coord()
	s := l.Score()
	return fmt.Sprintf("loc{%1.1f, %1.1f, %d}", l1, l2, s)
}

func (l *Loc) Score() uint64 {
	latNorm := round(normLat(l.lat))
	longNorm := round(normLong(l.long))

	score := interleave(latNorm, longNorm)
	return score
}

func (l *Loc) Coord() (grad, grad) {
	y := l.score >> 1
	x := l.score

	gridLat := compactInt64ToInt32(x)
	gridLong := compactInt64ToInt32(y)

	lat, long := convertGridNumbersToCoordinates(gridLat, gridLong)
	return grad(lat), grad(long)
}

func (l *Loc) Validate() bool {
	return ValidateCoord(l.long.v(), l.lat.v())
}

func Distance(l1, l2 Loc) float64 {
	lat1, lat2 := l1.lat.toRad(), l2.lat.toRad()
	dLat := lat2 - lat1
	dLon := l2.long.toRad() - l1.long.toRad()

	a := haversine(dLat) + math.Cos(lat1.v())*math.Cos(lat2.v())*haversine(dLon)
	// Floating-point error can push a slightly above 1; clamp, don't mask.
	a = math.Min(1, math.Max(0, a))
	return 2 * rEarth * math.Asin(math.Sqrt(a))
}

func haversine(a rad) float64 {
	return .5 * (1 - math.Cos(a.v()))
}

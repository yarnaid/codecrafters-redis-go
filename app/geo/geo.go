package geo

import (
	"fmt"
	"log/slog"
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
	lat1, long1 := l1.Coord()
	lat2, long2 := l2.Coord()
	dLat := lat2.toRad() - lat1.toRad()
	dLon := long2.toRad() - long1.toRad()
	slog.Debug("[Distance]", "dlat", dLat, "dLon", dLon)
	return 2 * rEarth * math.Asin(math.Sqrt(haversine(dLat)+
		math.Cos(l1.lat.toRad().v())*math.Cos(l2.lat.toRad().v())*haversine(dLon)))
}

func haversine(a rad) float64 {
	return .5 * (1 - math.Cos(a.v()))
}

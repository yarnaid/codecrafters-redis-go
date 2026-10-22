package geo

const (
	MinLong   = -180.
	MaxLong   = 180.
	MinLat    = -85.05112878
	MaxLat    = 85.05112878
	LatRange  = MaxLat - MinLat
	LongRange = MaxLong - MinLong
)

type Loc struct {
	lat, long float64
	score     uint64
}

func NewLoc(lat, long float64, score uint64) *Loc {
	return &Loc{lat: lat, long: long, score: score}
}

func (l *Loc) Score() uint64 {
	latNorm := round(normLat(l.lat))
	longNorm := round(normLong(l.long))

	score := interleave(latNorm, longNorm)
	return score
}

func (l *Loc) Coord() (float64, float64) {
	y := l.score >> 1
	x := l.score

	gridLat := compactInt64ToInt32(x)
	gridLong := compactInt64ToInt32(y)

	return convertGridNumbersToCoordinates(gridLat, gridLong)
}

func (l *Loc) Validate() bool {
	return ValidateCoord(l.long, l.lat)
}

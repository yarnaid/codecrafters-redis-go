package storage

import (
	"log/slog"

	"my-redis/app/geo"
	"my-redis/app/zset"
)

func (c *Coordinator) GeoDist(key, locName1, locName2 string) (float64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.geoDistLocked(key, locName1, locName2)
}

func (c *Coordinator) geoDistLocked(key, locName1, locName2 string) (float64, error) {
	score1, err := c.zScoreLocked(key, locName1)
	if err != nil {
		return 0, err
	}
	score2, err := c.zScoreLocked(key, locName2)
	if err != nil {
		return 0, err
	}
	loc1 := geo.NewLoc(0, 0, uint64(score1))
	loc2 := geo.NewLoc(0, 0, uint64(score2))
	return geo.Distance(*loc1, *loc2), nil
}

func (c *Coordinator) GeoSearch(key string, lon, lat, rad float64) ([]string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	interest := geo.NewLoc(lat, lon, 0)
	slog.Debug("[GeoSearch]", "center", interest, "lat", lat, "lon", lon)
	res := []string{}
	v, ok := c.backend.Get(key)
	if !ok {
		return nil, &NotFoundError{}
	}
	set, ok := v.Value.(zset.ZSet)
	if !ok {
		return nil, &WrongTypeError{v.Value, zset.NewZSetMapSlice()}
	}

	for name, score := range set.Values() {
		loc := geo.NewLoc(0, 0, uint64(score))
		d := geo.Distance(*interest, *loc)
		// slog.Debug("[GeoSearch]", "loc", loc, "dist", d, "rad", rad)
		if d <= rad {
			res = append(res, name)
		}
	}
	return res, nil
}

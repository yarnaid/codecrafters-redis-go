package storage

import (
	"log/slog"

	"my-redis/app/geo"
)

func (c *Coordinator) GeoDist(key, locName1, locName2 string) (float64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	score1, err := c.zScoreLocked(key, locName1)
	if err != nil {
		return 0, err
	}
	score2, err := c.zScoreLocked(key, locName2)
	if err != nil {
		return 0, err
	}
	loc1 := geo.NewLoc(0, 0, uint64(score1))
	loc1.Validate()
	loc2 := geo.NewLoc(0, 0, uint64(score2))
	loc2.Validate()
	res := geo.Distance(*loc1, *loc2)
	slog.Debug("[GeoDist]", "l1", loc1.String(), "l2", loc2.String(), "res", res)
	return res, nil
}

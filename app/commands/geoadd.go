package commands

import (
	"context"
	"fmt"
	"strconv"

	"my-redis/app/geo"
	"my-redis/app/parser"
)

type GeoAddCommand struct {
	Key       string
	Long, Lat float64
	Name      string
}

var _ Command = (*GeoAddCommand)(nil)

func (g *GeoAddCommand) Execute(ctx context.Context, env *Env, sess *Session) (parser.Serializable, error) {
	loc := geo.NewLoc(g.Lat, g.Long, 0)
	if !loc.Validate() {
		return nil, fmt.Errorf("ERR invalid longitude,latitude pair %f,%f", g.Long, g.Lat)
	}
	env.Coordinator.ZAdd(g.Key, g.Name, float64(loc.Score()))
	return parser.Int(1), nil
}

func parseGeoAdd(args []string) (Command, error) {
	long, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		return nil, err
	}
	lat, err := strconv.ParseFloat(args[2], 64)
	if err != nil {
		return nil, err
	}
	return &GeoAddCommand{args[0], long, lat, args[3]}, nil
}

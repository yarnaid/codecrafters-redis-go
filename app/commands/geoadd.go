package commands

import (
	"context"
	"fmt"
	"strconv"

	"my-redis/app/parser"
)

const (
	MinLong = -180.
	MaxLong = 180.
	MinLat  = -85.05112878
	MaxLat  = 85.05112878
)

type GeoAddCommand struct {
	Key       string
	Long, Lat float64
	Name      string
}

var _ Command = (*GeoAddCommand)(nil)

func (g *GeoAddCommand) Execute(ctx context.Context, env *Env, sess *Session) (parser.Serializable, error) {
	if !g.validateCoord() {
		return nil, fmt.Errorf("ERR invalid longitude,latitude pair %f,%f", g.Long, g.Lat)
	}
	env.Coordinator.ZAdd(g.Key, g.Name, 0)
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

func (g *GeoAddCommand) validateCoord() bool {
	return g.Long <= MaxLong && g.Long >= MinLong && g.Lat >= MinLat && g.Lat <= MaxLat
}

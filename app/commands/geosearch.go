package commands

import (
	"context"
	"strconv"

	"my-redis/app/parser"
)

type GeoSearchCommand struct {
	Key            string
	Long, Lat, Rad float64
}

var _ Command = (*GeoSearchCommand)(nil)

func (g *GeoSearchCommand) Execute(ctx context.Context, env *Env, sess *Session) (parser.Serializable, error) {
	res, err := env.Coordinator.GeoSearch(g.Key, g.Long, g.Lat, g.Rad)
	if err != nil {
		return nil, err
	}
	return parser.StringsToArray(res...), nil
}

func parseGeoSearch(args []string) (Command, error) {
	long, err := strconv.ParseFloat(args[2], 64)
	if err != nil {
		return nil, err
	}
	lat, err := strconv.ParseFloat(args[3], 64)
	if err != nil {
		return nil, err
	}
	rad, err := strconv.ParseFloat(args[5], 64)
	if err != nil {
		return nil, err
	}
	return &GeoSearchCommand{args[0], long, lat, rad}, nil
}

package commands

import (
	"context"
	"strconv"

	"my-redis/app/parser"
)

type GeoAddCommand struct {
	Long, Lat float64
	Name      string
}

var _ Command = (*GeoAddCommand)(nil)

func (g *GeoAddCommand) Execute(ctx context.Context, env *Env, sess *Session) (parser.Serializable, error) {
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
	return &GeoAddCommand{long, lat, args[3]}, nil
}

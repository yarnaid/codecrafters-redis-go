package commands

import (
	"context"
	"fmt"

	"my-redis/app/geo"
	"my-redis/app/parser"
)

type GeoPosCommand struct {
	Key   string
	Names []string
}

var _ Command = (*GeoPosCommand)(nil)

func (g *GeoPosCommand) Execute(ctx context.Context, env *Env, sess *Session) (parser.Serializable, error) {
	res := make(parser.Array[parser.Serializable], len(g.Names))
	for i, name := range g.Names {

		score, err := env.Coordinator.ZScore(g.Key, name)
		if err != nil {
			res[i] = parser.NullArray(0)
			continue
		}
		loc := geo.NewLoc(0, 0, uint64(score))
		lat, lon := loc.Coord()
		res[i] = parser.StringsToArray(fmt.Sprint(lon), fmt.Sprint(lat))
	}
	return res, nil
}

func parseGeoPos(args []string) (Command, error) {
	return &GeoPosCommand{args[0], args[1:]}, nil
}

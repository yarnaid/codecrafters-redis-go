package server

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

var logger = slog.New(tint.NewTextHandler(os.Stderr, &tint.Options{
	Level:      slog.LevelDebug,
	TimeFormat: "15:04:05",
}))

func init() {
	slog.SetDefault(logger)
}

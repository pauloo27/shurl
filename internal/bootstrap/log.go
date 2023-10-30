package bootstrap

import (
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
	"github.com/pauloo27/shurl/internal/config"
)

func setupLog(cfg *config.Config) {
	var handler slog.Handler
	level := cfg.Log.Level

	switch cfg.Log.Type {
	case config.LogTypeJSON:
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
		})
	case config.LogTypeColored:
		handler = tint.NewHandler(os.Stdout, &tint.Options{
			Level:      level,
			TimeFormat: time.DateTime,
		})
	default:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
		})
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
}

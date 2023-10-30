package bootstrap

import (
	"log/slog"
	"os"

	"github.com/pauloo27/shurl/internal/config"
)

const (
	DefaultConfigPath = "config.yaml"
)

func Start() {
	cfg, err := config.LoadConfig(DefaultConfigPath)
	if err != nil {
		slog.Error("Failed to load config:", "err", err)
		os.Exit(1)
	}

	setupLog(cfg)

	slog.Info("Starting shurl!")
	slog.Debug("If you can see this, debug logging is enabled!", "cool", true)
}

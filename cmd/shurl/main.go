package main

import (
	"log/slog"
	"os"

	"github.com/pauloo27/shurl/internal/bootstrap"
	"github.com/pauloo27/shurl/internal/config"
)

const (
	DefaultConfigPath = "config.yaml"
)

func main() {
	cfg, err := config.LoadConfigFromFile(DefaultConfigPath)
	if err != nil {
		slog.Error("Failed to load config:", "err", err)
		os.Exit(1)
	}

	bootstrap.Start(cfg)
}

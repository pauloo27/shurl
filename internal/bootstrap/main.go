package bootstrap

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
	"github.com/pauloo27/shurl/internal/config"
	"github.com/pauloo27/shurl/internal/ctx"
	"github.com/pauloo27/shurl/internal/providers/redis"
	"github.com/pauloo27/shurl/internal/server"
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

	rdb, err := redis.New(cfg.Redis)
	if err != nil {
		slog.Error("Failed to connect to redis:", tint.Err(err))
		os.Exit(1)
	}

	services := &ctx.Services{
		Config: cfg,
		Rdb:    rdb,
	}

	err = server.StartServer(services)
	if err != nil {
		slog.Error("Failed to start server:", tint.Err(err))
		os.Exit(1)
	}
}

package bootstrap

import (
	"log/slog"
	"os"

	"github.com/pauloo27/shurl/internal/config"
	"github.com/pauloo27/shurl/internal/ctx"
	"github.com/pauloo27/shurl/internal/providers/redis"
	"github.com/pauloo27/shurl/internal/server"
)

func Start(cfg *config.Config) {
	setupLog(cfg)

	slog.Info("Starting shurl!")
	slog.Debug("If you can see this, debug logging is enabled!", "cool", true)

	rdb, err := redis.New(cfg.Redis)
	if err != nil {
		slog.Error("Failed to connect to redis:", "err", err)
		os.Exit(1)
	}

	services := &ctx.Services{
		Config: cfg,
		Rdb:    rdb,
	}

	err = server.StartServer(services)
	if err != nil {
		slog.Error("Failed to start server:", "err", err)
		os.Exit(1)
	}
}

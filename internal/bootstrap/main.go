package bootstrap

import (
	"log/slog"
	"os"

	"github.com/pauloo27/shurl/internal/config"
	"github.com/pauloo27/shurl/internal/providers"
	"github.com/pauloo27/shurl/internal/providers/valkey"
	"github.com/pauloo27/shurl/internal/server"
)

func Start(cfg *config.Config) {
	setupLog(cfg)

	slog.Info("Starting shurl!")
	slog.Debug("If you can see this, debug logging is enabled!", "cool", true)

	vkey, err := valkey.New(cfg.Valkey)
	if err != nil {
		slog.Error("Failed to connect to valkey:", "err", err)
		os.Exit(1)
	}

	providers := &providers.Providers{
		Config: cfg,
		Valkey: vkey,
	}

	err = server.StartServer(providers)
	if err != nil {
		slog.Error("Failed to start server:", "err", err)
		os.Exit(1)
	}
}

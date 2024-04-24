package redict

import (
	"context"
	"log/slog"

	"github.com/pauloo27/shurl/internal/config"
	"github.com/redis/go-redis/v9"
)

func New(cfg *config.RedictConfig) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	cmd := rdb.Ping(context.Background())
	if err := cmd.Err(); err != nil {
		return nil, err
	}

	slog.Info("Connected to Redict", "addr", cfg.Address)

	return rdb, nil
}

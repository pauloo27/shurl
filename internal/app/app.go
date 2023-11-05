package app

import (
	"github.com/pauloo27/shurl/internal/config"
	"github.com/redis/go-redis/v9"
)

type App struct {
	Config *config.Config
	Rdb    *redis.Client

	Version string
}

func New(
	config *config.Config,
	rdb *redis.Client,
) *App {
	return &App{
		Config:  config,
		Rdb:     rdb,
		Version: "v0.0.1",
	}
}

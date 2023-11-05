package app

import "github.com/pauloo27/shurl/internal/config"

type App struct {
	Config *config.Config
}

func New(config *config.Config) *App {
	return &App{
		Config: config,
	}
}

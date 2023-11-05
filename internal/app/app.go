package app

import "github.com/pauloo27/shurl/internal/config"

type App struct {
	Config  *config.Config
	Version string
}

func New(config *config.Config) *App {
	return &App{
		Config:  config,
		Version: "v0.0.1",
	}
}

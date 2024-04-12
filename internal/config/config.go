package config

import (
	"log/slog"
)

type Config struct {
	Log    *LogConfig
	HTTP   *HTTPConfig
	Redict *RedictConfig

	Public *AppConfig

	Apps map[string]*AppConfig

	AppByAPIKey map[string]*AppConfig `yaml:"-" json:"-"`
}

type LogType string

const (
	LogTypeText    LogType = "text"
	LogTypeJSON    LogType = "json"
	LogTypeColored LogType = "colored"
)

type LogConfig struct {
	Level      slog.Level
	Type       LogType
	ShowSource bool
}

type HTTPConfig struct {
	Port int
}

type RedictConfig struct {
	Address  string
	Password string
	DB       int
}

type AppConfig struct {
	Enabled        bool
	APIKey         string
	AllowedDomains []string
	//LimitPerIPPerHour int TODO:
	//AllowCustomSlug bool TODO:
	//MinDurationSec  int TODO:
	//MaxDurationSec  int TODO:
}

package config

import (
	"log/slog"
)

type Config struct {
	Log    *LogConfig
	HTTP   *HTTPConfig
	Valkey *Valkey

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

type Valkey struct {
	Address  string
	Password string
	DB       int
}

type AppConfig struct {
	Enabled        bool
	APIKey         string
	MinDurationSec int
	MaxDurationSec int
	//LimitPerIPPerHour int TODO:
	//AllowCustomSlug bool TODO:
}

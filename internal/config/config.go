package config

import (
	"log/slog"
)

type Config struct {
	Log   *LogConfig
	HTTP  *HTTPConfig
	Redis *RedisConfig

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

type RedisConfig struct {
	Address  string
	Password string
	DB       int
}

type AppConfig struct {
	Enabled           bool
	APIKey            string
	LimitPerIPPerHour int
	AllowCustomSlug   bool
	AllowedDomains    []string
	MinDurationSec    int
	MaxDurationSec    int
}

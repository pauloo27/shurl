package config

import "log/slog"

type Config struct {
	Log   *LogConfig
	HTTP  *HTTPConfig
	Redis *RedisConfig
}

type LogType string

const (
	LogTypeText    LogType = "text"
	LogTypeJSON    LogType = "json"
	LogTypeColored LogType = "colored"
)

type LogConfig struct {
	Level slog.Level
	Type  LogType
}

type HTTPConfig struct {
	Port int
}

type RedisConfig struct {
	Address  string
	Password string
	DB       int
}

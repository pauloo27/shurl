package config

import "log/slog"

type Config struct {
	Log  LogConfig
	Http HttpConfig
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

type HttpConfig struct {
	Port int
}

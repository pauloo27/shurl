package config

import (
	"errors"
	"os"

	"github.com/ghodss/yaml"
)

func LoadConfigFromFile(configPath string) (*Config, error) {
	/* #nosec G304 */
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	return LoadConfigFromData(data)
}

func LoadConfigFromData(data []byte) (*Config, error) {
	var config Config

	err := yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	ensureNotNil(&config)

	config.AppByAPIKey = make(map[string]*AppConfig)

	if config.Public.APIKey != "" {
		return nil, errors.New("public client must not have api key")
	}

	for _, app := range config.Apps {
		config.AppByAPIKey[app.APIKey] = app
	}

	return &config, nil
}

func ensureNotNil(cfg *Config) {
	if cfg.Log == nil {
		cfg.Log = &LogConfig{}
	}
	if cfg.HTTP == nil {
		cfg.HTTP = &HTTPConfig{}
	}
	if cfg.Valkey == nil {
		cfg.Valkey = &Valkey{}
	}
	if cfg.Public == nil {
		cfg.Public = &AppConfig{}
	}
	if cfg.Apps == nil {
		cfg.Apps = make(map[string]*AppConfig)
	}
}

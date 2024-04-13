package config

import (
	"errors"
	"os"

	"github.com/ghodss/yaml"
)

func LoadConfig(configPath string) (*Config, error) {
	/* #nosec G304 */
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	config.AppByAPIKey = make(map[string]*AppConfig)

	if config.Public.APIKey != "" {
		return nil, errors.New("public client must not have api key")
	}

	for _, app := range config.Apps {
		config.AppByAPIKey[app.APIKey] = app
	}

	return &config, nil
}

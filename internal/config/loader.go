package config

import (
	"github.com/ghodss/yaml"
	"os"
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

	return &config, nil
}

package mocker

import "github.com/pauloo27/shurl/internal/config"

func MakeConfigMock(cfg *config.Config) *config.Config {
	if cfg == nil {
		return &config.Config{}
	}

	cfg.AppByAPIKey = make(map[string]*config.AppConfig)
	for _, app := range cfg.Apps {
		cfg.AppByAPIKey[app.APIKey] = app
	}

	return cfg
}

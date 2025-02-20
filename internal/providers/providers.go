package providers

import (
	"github.com/pauloo27/shurl/internal/config"
	"github.com/valkey-io/valkey-go"
)

type Providers struct {
	Config *config.Config
	Valkey valkey.Client
}

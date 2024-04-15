package ctx

import (
	"context"

	"github.com/pauloo27/shurl/internal/config"
	"github.com/redis/go-redis/v9"
)

type Key string

const (
	ProvidersKey Key = "providers"
)

type Providers struct {
	Config *config.Config
	Rdb    *redis.Client
}

func GetProviders(ctx context.Context) *Providers {
	return ctx.Value(ProvidersKey).(*Providers)
}

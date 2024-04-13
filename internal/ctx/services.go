package ctx

import (
	"context"

	"github.com/pauloo27/shurl/internal/config"
	"github.com/redis/go-redis/v9"
)

type Key string

const (
	ServicesKey Key = "services"
)

type Services struct {
	Config *config.Config
	Rdb    *redis.Client
}

func GetServices(ctx context.Context) *Services {
	return ctx.Value(ServicesKey).(*Services)
}

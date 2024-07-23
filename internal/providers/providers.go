package providers

import (
	"github.com/pauloo27/shurl/internal/config"
	"github.com/redis/go-redis/v9"
)

type Providers struct {
	Config *config.Config
	Rdb    *redis.Client
}

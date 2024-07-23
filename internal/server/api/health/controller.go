package health

import (
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

type HealthController struct {
	rdb *redis.Client
}

func NewHealthController(rdb *redis.Client) *HealthController {
	return &HealthController{rdb}
}

func (c *HealthController) Route(e *echo.Echo) {
	e.GET("/api/v1/healthz", c.Health)
}

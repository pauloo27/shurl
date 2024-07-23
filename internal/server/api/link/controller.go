package link

import (
	"github.com/labstack/echo/v4"
	"github.com/pauloo27/shurl/internal/config"
	"github.com/redis/go-redis/v9"
)

type LinkController struct {
	rdb *redis.Client
	cfg *config.Config
}

func NewLinkController(cfg *config.Config, rdb *redis.Client) *LinkController {
	return &LinkController{rdb, cfg}
}

func (c *LinkController) Route(e *echo.Echo) {
	e.POST("/api/v1/links", c.Create)
	e.GET("/:slug", c.Redirect)
}

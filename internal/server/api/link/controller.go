package link

import (
	"github.com/labstack/echo/v4"
	"github.com/pauloo27/shurl/internal/config"
	"github.com/valkey-io/valkey-go"
)

type LinkController struct {
	vkey valkey.Client
	cfg  *config.Config
}

func NewLinkController(cfg *config.Config, vkey valkey.Client) *LinkController {
	return &LinkController{vkey, cfg}
}

func (c *LinkController) Route(e *echo.Echo) {
	e.POST("/api/v1/links", c.Create)
	e.GET("/:slug", c.Redirect)
}

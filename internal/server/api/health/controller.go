package health

import (
	"github.com/labstack/echo/v4"
	"github.com/valkey-io/valkey-go"
)

type HealthController struct {
	vkey valkey.Client
}

func NewHealthController(vkey valkey.Client) *HealthController {
	return &HealthController{vkey}
}

func (c *HealthController) Route(e *echo.Echo) {
	e.GET("/api/v1/healthz", c.Health)
}

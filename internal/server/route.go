package server

import (
	"github.com/labstack/echo/v4"
	"github.com/pauloo27/shurl/internal/providers"
	"github.com/pauloo27/shurl/internal/server/api/health"
	"github.com/pauloo27/shurl/internal/server/api/link"

	// swagger :D
	_ "github.com/pauloo27/shurl/internal/server/docs"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func route(providers *providers.Providers, e *echo.Echo) {
	e.GET("/", func(c echo.Context) error {
		return c.Redirect(301, "/api/v1/swagger/index.html")
	})

	routeHealth(providers, e)
	routeLink(providers, e)
	routeSwagger(e)
}

func routeSwagger(g *echo.Echo) {
	g.GET("/api/v1/swagger/*", echoSwagger.WrapHandler)
}

func routeHealth(providers *providers.Providers, e *echo.Echo) {
	c := health.NewHealthController(providers.Valkey)
	c.Route(e)
}

func routeLink(providers *providers.Providers, e *echo.Echo) {
	c := link.NewLinkController(providers.Config, providers.Valkey)
	c.Route(e)
}

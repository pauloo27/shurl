package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/pauloo27/shurl/internal/providers"
)

// @title			Shurl API
// @version		1.0
// @description	URL Shortener API
// @license.name	MIT
// @license.url	https://opensource.org/licenses/MIT
// @BasePath		/api/v1
func StartServer(providers *providers.Providers) error {
	e := echo.New()

	bindAddr := fmt.Sprintf(":%d", providers.Config.HTTP.Port)

	route(providers, e)

	server := &http.Server{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Addr:         bindAddr,
	}

	return e.StartServer(server)
}

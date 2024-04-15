package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chi_middleware "github.com/go-chi/chi/v5/middleware"
	"github.com/pauloo27/shurl/internal/ctx"
	"github.com/pauloo27/shurl/internal/server/middleware"
	"github.com/pauloo27/shurl/internal/server/router"
)

// @title			Shurl API
// @version		1.0
// @description	URL Shortener API
// @license.name	MIT
// @license.url	https://opensource.org/licenses/MIT
// @BasePath		/api/v1
func StartServer(services *ctx.Services) error {
	r := chi.NewRouter()

	r.Use(chi_middleware.RequestID)
	r.Use(chi_middleware.RealIP)
	r.Use(chi_middleware.Recoverer)

	r.Use(middleware.LoggerMiddleware)
	r.Use(middleware.ServicesContext(services))

	router.RouteApp(r)

	bindAddr := fmt.Sprintf(":%d", services.Config.HTTP.Port)

	slog.Info("Starting server", "addr", bindAddr)

	server := &http.Server{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      r,
		Addr:         bindAddr,
	}

	return server.ListenAndServe()
}

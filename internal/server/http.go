package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/pauloo27/shurl/internal/ctx"
	"github.com/pauloo27/shurl/internal/server/router"
)

func StartServer(services *ctx.Services) error {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(loggerMiddleware)
	r.Use(middleware.Recoverer)
	r.Use(servicesContext(services))

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

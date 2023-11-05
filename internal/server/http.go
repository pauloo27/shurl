package server

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/pauloo27/shurl/internal/app"
)

func StartServer(shurl *app.App) error {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(loggerMiddleware)
	r.Use(middleware.Recoverer)

	routeApp(shurl, r)

	bindAddr := fmt.Sprintf(":%d", shurl.Config.Http.Port)

	slog.Info("Starting server", "addr", bindAddr)
	return http.ListenAndServe(bindAddr, r)
}

package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/pauloo27/shurl/internal/server/api/health"
	"github.com/pauloo27/shurl/internal/server/api/link"
)

func RouteApp(root *chi.Mux) {
	root.Route("/api/v1", func(r chi.Router) {
		r.Get("/healthz", health.Health)
		r.Post("/links", link.Create)
	})

	root.Get("/{slug}", link.Redirect)
}

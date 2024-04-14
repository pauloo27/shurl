package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/pauloo27/shurl/internal/server/api/health"
	"github.com/pauloo27/shurl/internal/server/api/link"
	"github.com/swaggo/http-swagger/v2"
	// swagger :)
	_ "github.com/pauloo27/shurl/internal/server/docs"
)

func RouteApp(root *chi.Mux) {
	root.Route("/api/v1", func(r chi.Router) {
		r.Get("/healthz", health.Health)
		r.Post("/links", link.Create)
		r.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL("/api/v1/swagger/doc.json"),
		))
	})

	root.Get("/{slug}", link.Redirect)
}

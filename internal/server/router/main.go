package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/pauloo27/shurl/internal/server/api/health"
	"github.com/pauloo27/shurl/internal/server/api/link"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	// swagger :)
	_ "github.com/pauloo27/shurl/internal/server/docs"
)

func RouteApp(root *chi.Mux) {
	root.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/api/v1/swagger/index.html", http.StatusSeeOther)
	})

	root.Route("/api/v1", func(r chi.Router) {
		r.Get("/healthz", wrap(health.Health))
		r.Post("/links", wrap(link.Create))
		r.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL("/api/v1/swagger/doc.json"),
		))
		r.Get("/swagger", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/api/v1/swagger/index.html", http.StatusSeeOther)
		})
	})

	root.Get("/{slug}", wrap(link.Redirect))
}

package router

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/pauloo27/shurl/internal/ctx"
	"github.com/pauloo27/shurl/internal/server/api/link"
)

func RouteApp(services *ctx.Services, root *chi.Mux) {
	root.Route("/_", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Running shurl %s", services.Version)
		})
	})

	root.Get("/{slug}", link.Redirect)
}

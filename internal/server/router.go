package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/pauloo27/shurl/internal/app"
)

func routeApp(shurl *app.App, root *chi.Mux) {
	root.Route("/_", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Running shurl %s", shurl.Version)
		})
	})
}

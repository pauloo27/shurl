package middleware

import (
	"context"
	"net/http"

	"github.com/pauloo27/shurl/internal/ctx"
)

func ServicesContext(services *ctx.Services) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(
				w,
				r.WithContext(
					context.WithValue(r.Context(), ctx.ServicesKey, services),
				),
			)
		})
	}
}

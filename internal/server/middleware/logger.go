package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/pauloo27/shurl/internal/ctx"
)

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		start := time.Now()

		id := r.Context().Value(middleware.RequestIDKey)
		providers := ctx.GetProviders(r.Context())

		log := slog.With("id", id)

		log.Info(
			"Http request",
			"remote_addr", r.RemoteAddr,
			"method", r.Method,
			"url", r.URL,
		)

		providers.Logger = log

		next.ServeHTTP(ww, r)
		status := ww.Status()

		log.Info(
			"Http response",
			"status", status,
			"took", time.Since(start),
		)
	})
}

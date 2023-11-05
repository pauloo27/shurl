package link

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/lmittmann/tint"
	"github.com/pauloo27/shurl/internal/server/ctx"
	"github.com/redis/go-redis/v9"
)

func Redirect(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	rdb := ctx.GetApp(c).Rdb

	slug := chi.URLParam(r, "slug")

	res := rdb.Get(c, fmt.Sprintf("link:%s", slug))
	if err := res.Err(); err != nil {
		if errors.Is(err, redis.Nil) {
			slog.Warn("Link not found", "slug", slug)
			http.NotFound(w, r)
			return
		}
		slog.Error("Failed to get link", "slug", slug, "err", tint.Err(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	slog.Info("Redirecting", "slug", slug, "url", res.Val())

	http.Redirect(w, r, res.Val(), http.StatusTemporaryRedirect)
}

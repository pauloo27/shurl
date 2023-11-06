package link

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/lmittmann/tint"
	"github.com/pauloo27/shurl/internal/ctx"
	"github.com/pauloo27/shurl/internal/server/api"
	"github.com/redis/go-redis/v9"
)

func Redirect(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	rdb := ctx.GetServices(c).Rdb

	slug := chi.URLParam(r, "slug")

	res := rdb.Get(c, fmt.Sprintf("link:%s", slug))
	if err := res.Err(); err != nil {
		if errors.Is(err, redis.Nil) {
			slog.Warn("Link not found", "slug", slug)
			api.Err(w, http.StatusNotFound, api.NotFoundErr, "Link not found")
			return
		}
		slog.Error("Failed to get link", "slug", slug, tint.Err(err))
		api.Err(w, http.StatusInternalServerError, api.InternalServerErr, "Something went wrong")
		return
	}

	slog.Info("Redirecting", "slug", slug, "url", res.Val())

	http.Redirect(w, r, res.Val(), http.StatusTemporaryRedirect)
}

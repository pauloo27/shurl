package link

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/pauloo27/shurl/internal/ctx"
	"github.com/pauloo27/shurl/internal/server/api"
	"github.com/redis/go-redis/v9"
)

func Redirect(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	rdb := ctx.GetServices(c).Rdb

	domain := r.Host
	slug := chi.URLParam(r, "slug")

	key := fmt.Sprintf("link:%s/%s", domain, slug)

	res := rdb.Get(c, key)

	if err := res.Err(); err != nil {
		if errors.Is(err, redis.Nil) {
			slog.Warn("Link not found", "domain", domain, "slug", slug)
			api.Err(w, http.StatusNotFound, api.NotFoundErr, "Link not found")
			return
		}
		slog.Error("Failed to get link", "slug", slug, "err", err)
		api.Err(w, http.StatusInternalServerError, api.InternalServerErr, "Something went wrong")
		return
	}

	slog.Info("Redirecting", "domain", domain, "slug", slug, "url", res.Val())

	http.Redirect(w, r, res.Val(), http.StatusTemporaryRedirect)
}

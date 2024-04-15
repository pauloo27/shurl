package link

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/pauloo27/shurl/internal/ctx"
	"github.com/pauloo27/shurl/internal/server/api"
	"github.com/redis/go-redis/v9"
)

// Redirect godoc
//
//	@Summary		Redirect to the original URL
//	@Description	Redirect from domain/slug to the original URL
//	@Tags			link
//	@Param			slug	path	string	true	"Slug to redirect from"
//	@Success		307
//	@Failure		404	{object}	api.Error[map[string]string]	"Link not found"
//	@Failure		500	{object}	api.Error[map[string]string]	"Internal server error"
//	@Router			/{slug} [get]
func Redirect(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	rdb := ctx.GetServices(c).Rdb

	domain := r.Host
	slug := r.PathValue("slug")

	key := fmt.Sprintf("link:%s/%s", domain, slug)

	res := rdb.Get(c, key)

	if err := res.Err(); err != nil {
		if errors.Is(err, redis.Nil) {
			slog.Warn("Link not found", "domain", domain, "slug", slug)
			api.Err(w, api.NotFoundErr, "Link not found")
			return
		}
		slog.Error("Failed to get link", "slug", slug, "err", err)
		api.Err(w, api.InternalServerErr, "Something went wrong")
		return
	}

	slog.Info("Redirecting", "domain", domain, "slug", slug, "url", res.Val())

	http.Redirect(w, r, res.Val(), http.StatusTemporaryRedirect)
}

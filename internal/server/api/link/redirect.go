package link

import (
	"errors"
	"fmt"
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
func Redirect(r *http.Request) api.Response {
	c := r.Context()
	providers := ctx.GetProviders(c)
	rdb := providers.Rdb
	log := providers.Logger

	domain := r.Host
	slug := r.PathValue("slug")

	key := fmt.Sprintf("link:%s/%s", domain, slug)

	res := rdb.Get(c, key)

	if err := res.Err(); err != nil {
		if errors.Is(err, redis.Nil) {
			return api.Err(api.NotFoundErr, "Link not found")
		}
		log.Error("Failed to get link", "slug", slug, "err", err)
		return api.Err(api.InternalServerErr, "Something went wrong")
	}

	return api.Redirect(res.Val())
}

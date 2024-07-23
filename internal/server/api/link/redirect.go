package link

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
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
//	@Failure		404	{object}	api.NotFoundError		"Link not found"
//	@Failure		500	{object}	api.InternalServerError	"Internal server error"
//	@Router			/{slug} [get]
func (c *LinkController) Redirect(ctx echo.Context) error {
	domain := ctx.Request().Host
	slug := ctx.Param("slug")

	slog.Info("h-hello?", "slug", slug, "domain", domain)
	key := fmt.Sprintf("link:%s/%s", domain, slug)

	res := c.rdb.Get(context.Background(), key)

	if err := res.Err(); err != nil {
		if errors.Is(err, redis.Nil) {
			return ctx.JSON(api.Err(api.ErrNotFound, "Link not found"))
		}
		slog.Error("Failed to get link", "slug", slug, "err", err)
		return ctx.JSON(api.Err(api.ErrInternalServer, "Something went wrong"))
	}

	return ctx.Redirect(http.StatusTemporaryRedirect, res.Val())
}

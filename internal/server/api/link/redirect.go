package link

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pauloo27/shurl/internal/server/api"
	"github.com/valkey-io/valkey-go"
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

	cmd := c.vkey.B().Get().Key(key).Build()
	res := c.vkey.Do(context.Background(), cmd)

	if err := res.Error(); err != nil {
		if valkey.IsValkeyNil(err) {
			return ctx.JSON(api.Err(api.ErrNotFound, "Link not found"))
		}
		slog.Error("Failed to get link", "slug", slug, "err", err)
		return ctx.JSON(api.Err(api.ErrInternalServer, "Something went wrong"))
	}

	value, err := res.ToString()
	if err != nil {
		slog.Error("Failed to parse string value", "slug", slug, "err", err)
		return ctx.JSON(api.Err(api.ErrInternalServer, "Something went wrong"))
	}

	return ctx.Redirect(http.StatusTemporaryRedirect, value)
}

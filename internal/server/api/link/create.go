package link

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/lmittmann/tint"
	"github.com/pauloo27/shurl/internal/ctx"
	"github.com/pauloo27/shurl/internal/server/api"
	"github.com/pauloo27/shurl/internal/server/validator"
)

type CreateLinkBody struct {
	Slug        string `json:"slug" validate:"required,min=3,max=20"`
	OriginalURL string `json:"original_url" validate:"required,http_url"`
}

func Create(w http.ResponseWriter, r *http.Request) {
	body, ok := validator.MustGetBody[CreateLinkBody](w, r)
	if !ok {
		return
	}

	c := r.Context()
	rdb := ctx.GetServices(c).Rdb

	slog.Info("Creating link", "slug", body.Slug, "url", body.OriginalURL)
	res := rdb.Set(c, fmt.Sprintf("link:%s", body.Slug), body.OriginalURL, 0)

	if err := res.Err(); err != nil {
		slog.Error("Failed to create link", "slug", body.Slug, tint.Err(err))
		api.Err(w, http.StatusInternalServerError, api.InternalServerErr, "Something went wrong")
		return
	}

	api.Created(w, body)
}

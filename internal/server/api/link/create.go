package link

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/lmittmann/tint"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/pauloo27/shurl/internal/ctx"
	"github.com/pauloo27/shurl/internal/server/api"
	"github.com/pauloo27/shurl/internal/server/validator"
)

type CreateLinkBody struct {
	Slug        string `json:"slug,omitempty" validate:"omitempty,min=3,max=20"`
	Domain      string `json:"domain,omitempty" validate:"omitempty,hostname"`
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

	slug := body.Slug
	if slug == "" {
		randomSlug, err := gonanoid.New()
		if err != nil {
			slog.Error("Failed to generate random slug", tint.Err(err))
			api.Err(w, http.StatusInternalServerError, api.InternalServerErr, "Something went wrong")
			return
		}
		slug = randomSlug
	}

	domain := body.Domain
	if domain == "" {
		// TODO: get from config
		domain = "localhost"
	}

	// TODO: create a model type
	link := map[string]any{
		"slug":       slug,
		"domain":     domain,
		"url":        body.OriginalURL,
		"created_at": time.Now(),
	}

	// TODO: check if slug already exists
	res := rdb.HSet(c, fmt.Sprintf("link:%s", slug), link)

	if err := res.Err(); err != nil {
		slog.Error("Failed to create link", "slug", body.Slug, tint.Err(err))
		api.Err(w, http.StatusInternalServerError, api.InternalServerErr, "Something went wrong")
		return
	}

	api.Created(w, link)
}

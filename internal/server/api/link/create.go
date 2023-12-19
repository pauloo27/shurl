package link

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/lmittmann/tint"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/pauloo27/shurl/internal/config"
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
	services := ctx.GetServices(c)
	rdb := services.Rdb
	cfg := services.Config

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

	var app *config.AppConfig

	apiKey := r.Header.Get("X-API-Key")
	if apiKey == "" {
		app = cfg.Public
	} else {
		app = cfg.AppByAPIKey[apiKey]
	}
	if app == nil || !app.Enabled {
		api.Err(w, http.StatusUnauthorized, api.UnauthorizedErr, "Invalid API key")
		return
	}

	domain := body.Domain
	if domain == "" {
		domain = app.AllowedDomains[0]
	}

	// TODO: create a model type
	link := map[string]any{
		"slug":       slug,
		"domain":     domain,
		"url":        body.OriginalURL,
		"created_at": time.Now(),
	}

	// TODO: use TX; check if slug already exists
	res := rdb.HSet(c, fmt.Sprintf("link:%s", slug), link)

	if err := res.Err(); err != nil {
		slog.Error("Failed to create link", "slug", body.Slug, tint.Err(err))
		api.Err(w, http.StatusInternalServerError, api.InternalServerErr, "Something went wrong")
		return
	}

	api.Created(w, link)
}

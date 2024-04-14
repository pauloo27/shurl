package link

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/pauloo27/shurl/internal/config"
	"github.com/pauloo27/shurl/internal/ctx"
	"github.com/pauloo27/shurl/internal/models"
	"github.com/pauloo27/shurl/internal/server/api"
	"github.com/pauloo27/shurl/internal/server/validator"
)

var (
	SlugBlacklist = map[string]bool{
		"api":        true,
		"links":      true,
		"admin":      true,
		"robots.txt": true,
	}
)

const (
	DefaultSlugLength = 8
)

type CreateLinkBody struct {
	Slug        string `json:"slug" validate:"omitempty,min=3,max=20"`
	Domain      string `json:"domain" validate:"omitempty,min=1"`
	OriginalURL string `json:"original_url" validate:"required,http_url"`
	TTL         *int   `json:"ttl" validate:"required"`
}

func Create(w http.ResponseWriter, r *http.Request) {
	body, ok := validator.MustGetBody[CreateLinkBody](w, r)
	if !ok {
		return
	}

	c := r.Context()
	services := ctx.GetServices(c)
	cfg := services.Config
	rdb := services.Rdb

	slug := body.Slug
	if slug == "" {
		randomSlug, err := gonanoid.New(DefaultSlugLength)
		if err != nil {
			slog.Error("Failed to generate random slug", "err", err)
			api.Err(w, http.StatusInternalServerError, api.InternalServerErr, "Something went wrong")
			return
		}
		slug = randomSlug
		slog.Info("Slug not provided, generating a random one", "slug", slug)
	}

	if SlugBlacklist[slug] {
		api.Err(w, http.StatusForbidden, api.ForbiddenErr, "Slug is blacklisted")
		return
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
		if len(app.AllowedDomains) == 0 {
			api.Err(w, http.StatusForbidden, api.ForbiddenErr, "No allowed domains for this app")
			slog.Info("Missing allowed domains for app", "apiKey", app.APIKey)
			return
		}
		slog.Info("Domain not provided, using the first allowed domain from app", "domain", domain)
		domain = app.AllowedDomains[0]
	} else {
		allowed := false
		for _, allowedDomain := range app.AllowedDomains {
			if domain == allowedDomain {
				allowed = true
				break
			}
		}
		if !allowed {
			api.Err(w, http.StatusForbidden, api.ForbiddenErr, "Domain not allowed for this app")
			return
		}
	}

	slog.Info("Creating link", "domain", domain, "slug", slug, "url", body.OriginalURL)

	ttlInSecs := *body.TTL
	var ttl time.Duration

	if ttlInSecs != 0 {
		ttl = time.Duration(ttlInSecs) * time.Second
	}

	link := models.Link{
		Slug:        slug,
		Domain:      domain,
		OriginalURL: body.OriginalURL,
		TTL:         ttlInSecs,
	}

	key := fmt.Sprintf("link:%s/%s", domain, slug)
	cmd := rdb.SetNX(context.Background(), key, body.OriginalURL, ttl)
	if cmd.Err() != nil {
		slog.Error("Failed to create link", "err", cmd.Err())
		api.Err(w, http.StatusInternalServerError, api.InternalServerErr, "Something went wrong")
		return
	}

	if !cmd.Val() {
		slog.Error("Link already exists", "slug", slug)
		api.Err(w, http.StatusConflict, api.ConflictErr, "Link already exists")
		return
	}

	api.Created(w, link)
}

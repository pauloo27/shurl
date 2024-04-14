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
	Slug        string `json:"slug" validate:"omitempty,min=3,max=20,excludes=/"`
	Domain      string `json:"domain" validate:"omitempty,min=1"`
	OriginalURL string `json:"original_url" validate:"required,http_url"`
	TTL         *int   `json:"ttl" validate:"required"`
}

// Create godoc
//
//	@Summary		Create a link
//	@Description	Create a link from a slug to the original URL.
//	@Description	If no slug is provided, a random one will be generated.
//	@Description	If no domain is provided, the first allowed domain from the app will be used.
//	@Description	The ttl is required. 0 means no expiration, otherwise it's the number of seconds until expiration.
//	@Param			body	body	CreateLinkBody	true	"Domain and slug are optional"
//	@Tags			link
//	@Produce		json
//	@Router			/links [post]
//	@Success		201	{object}	models.Link
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
			api.Err(w, api.InternalServerErr, "Something went wrong")
			return
		}
		slug = randomSlug
		slog.Info("Slug not provided, generating a random one", "slug", slug)
	}

	if SlugBlacklist[slug] {
		api.Err(w, api.ForbiddenErr, "Slug is blacklisted")
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
		api.Err(w, api.UnauthorizedErr, "Invalid API key")
		return
	}

	domain := body.Domain
	if domain == "" {
		if len(app.AllowedDomains) == 0 {
			api.Err(w, api.ForbiddenErr, "No allowed domains for this app")
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
			api.Err(w, api.ForbiddenErr, "Domain not allowed for this app")
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
		api.Err(w, api.InternalServerErr, "Something went wrong")
		return
	}

	if !cmd.Val() {
		slog.Error("Link already exists", "slug", slug)
		api.Err(w, api.ConflictErr, "Link already exists")
		return
	}

	api.Created(w, link)
}

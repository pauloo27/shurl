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
	DefaultSlugLength = 6
)

type CreateLinkBody struct {
	Slug        string `json:"slug" validate:"omitempty,min=3,max=20,excludes=/"`
	Domain      string `json:"domain" validate:"omitempty,min=1"`
	OriginalURL string `json:"original_url" validate:"required,http_url"`
	TTL         *int   `json:"ttl" validate:"min=0,max=31536000"`
}

// Create godoc
//
//	@Summary		Create a link
//	@Description	Create a link from a slug to the original URL.
//	@Description	If no slug is provided, a random one will be generated.
//	@Description	If no domain is provided, the first allowed domain from the app will be used.
//	@Description	The ttl is required. 0 means no expiration, otherwise it's the number of seconds until expiration.
//	@Description	The ttl can't be greater than 1 year (31536000 seconds).
//	@Description	The API Key may limit the allowed domains and the ttl.
//	@Param			body	body	CreateLinkBody	true	"Domain and slug are optional"
//	@Tags			link
//	@Produce		json
//	@Router			/links [post]
//	@Success		201	{object}	models.Link						"Created"
//	@Failure		400	{object}	api.Error[map[string]string]	"Bad request"
//	@Failure		500	{object}	api.Error[map[string]string]	"Internal server error"
//	@Failure		401	{object}	api.Error[map[string]string]	"Missing API Key"
//	@Failure		403	{object}	api.Error[map[string]string]	"Invalid API Key"
//	@Failure		409	{object}	api.Error[map[string]string]	"Duplicated link"
//	@Failure		422	{object}	api.Error[map[string]string]	"Validation error"
//	@Security		ApiKeyAuth
//	@Param			X-API-Key	header	string	false	"API Key, leave empty for public access (if enabled in the server)"
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

	if app.MaxDurationSec != 0 && ttlInSecs > app.MaxDurationSec {
		api.Err(w, api.BadRequestErr, fmt.Sprintf("TTL too high, max is %d", app.MaxDurationSec))
		return
	}

	if app.MinDurationSec != 0 && ttlInSecs < app.MinDurationSec {
		api.Err(w, api.BadRequestErr, fmt.Sprintf("TTL too low, min is %d", app.MinDurationSec))
		return
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

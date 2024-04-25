package link

import (
	"context"
	"fmt"
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
func Create(r *http.Request) api.Response {
	body, validationErr := validator.MustGetBody[CreateLinkBody](r)
	if validationErr != nil {
		return api.DetailedError(validationErr.Error, validationErr.Details)
	}

	c := r.Context()
	providers := ctx.GetProviders(c)
	cfg := providers.Config
	rdb := providers.Rdb
	log := providers.Logger

	slug := body.Slug
	if slug == "" {
		randomSlug, err := gonanoid.New(DefaultSlugLength)
		if err != nil {
			log.Error("Failed to generate random slug", "err", err)
			return api.Err(api.InternalServerErr, "Something went wrong")
		}
		slug = randomSlug
	}

	if SlugBlacklist[slug] {
		return api.Err(api.ForbiddenErr, "Slug is blacklisted")
	}

	var app *config.AppConfig

	apiKey := r.Header.Get("X-API-Key")
	if apiKey == "" {
		app = cfg.Public
	} else {
		app = cfg.AppByAPIKey[apiKey]
	}

	if app == nil || !app.Enabled {
		return api.Err(api.UnauthorizedErr, "Invalid API key")
	}

	domain := body.Domain
	if domain == "" {
		if len(app.AllowedDomains) == 0 {
			log.Info("Missing allowed domains for app", "apiKey", app.APIKey)
			return api.Err(api.ForbiddenErr, "No allowed domains for this app")
		}
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
			return api.Err(api.ForbiddenErr, "Domain not allowed for this app")
		}
	}

	log.Info("Creating link", "domain", domain, "slug", slug, "url", body.OriginalURL)

	ttlInSecs := *body.TTL
	var ttl time.Duration

	if ttlInSecs != 0 {
		ttl = time.Duration(ttlInSecs) * time.Second
	}

	if app.MaxDurationSec != 0 && ttlInSecs > app.MaxDurationSec {
		return api.Err(api.BadRequestErr, fmt.Sprintf("TTL too high, max is %d", app.MaxDurationSec))
	}

	if app.MinDurationSec != 0 && ttlInSecs < app.MinDurationSec {
		return api.Err(api.BadRequestErr, fmt.Sprintf("TTL too low, min is %d", app.MinDurationSec))
	}

	link := models.Link{
		Slug:        slug,
		Domain:      domain,
		OriginalURL: body.OriginalURL,
		TTL:         ttlInSecs,
		URL:         fmt.Sprintf("https://%s/%s", domain, slug),
	}

	key := fmt.Sprintf("link:%s/%s", domain, slug)
	cmd := rdb.SetNX(context.Background(), key, body.OriginalURL, ttl)

	if cmd.Err() != nil {
		return api.Err(api.InternalServerErr, "Something went wrong")
	}

	if !cmd.Val() {
		return api.Err(api.ConflictErr, "Link already exists")
	}

	return api.Created(link)
}

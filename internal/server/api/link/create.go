package link

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/pauloo27/shurl/internal/config"
	"github.com/pauloo27/shurl/internal/models"
	"github.com/pauloo27/shurl/internal/server/api"
	"github.com/pauloo27/shurl/internal/server/core/validator"
	"github.com/valkey-io/valkey-go"
)

var (
	slugBlacklist = map[string]bool{
		"api":        true,
		"links":      true,
		"admin":      true,
		"robots.txt": true,
	}
)

const (
	defaultSlugLength = 6
)

type CreateLinkBody struct {
	Slug        string `json:"slug" validate:"omitempty,min=3,max=20,excludes=/"`
	OriginalURL string `json:"original_url" validate:"required,http_url"`
	TTL         *int   `json:"ttl" validate:"min=0,max=31536000"`
}

// Create godoc
//
//	@Summary		Create a link
//	@Description	Create a link from a slug to the original URL.
//	@Description	If no slug is provided, a random one will be generated.
//	@Description	The ttl is required. 0 means no expiration, otherwise it's the number of seconds until expiration.
//	@Description	The ttl can't be greater than 1 year (31536000 seconds).
//	@Description	The API Key may limit the ttl.
//	@Param			body	body	CreateLinkBody	true	"Slug is optional"
//	@Tags			link
//	@Produce		json
//	@Router			/links [post]
//	@Success		201	{object}	models.Link				"Created"
//	@Failure		400	{object}	api.BadRequestError		"Bad request"
//	@Failure		500	{object}	api.InternalServerError	"Internal server error"
//	@Failure		401	{object}	api.UnauthorizedError	"Missing API Key"
//	@Failure		403	{object}	api.ForbiddenError		"Invalid API Key"
//	@Failure		409	{object}	api.ConflictError		"Duplicated link"
//	@Failure		422	{object}	api.ValidationError		"Validation error"
//	@Security		ApiKeyAuth
//	@Param			X-API-Key	header	string	false	"API Key, leave empty for public access (if enabled in the server)"
func (c *LinkController) Create(ctx echo.Context) error {
	body, validationErr := validator.MustBindAndValidate[CreateLinkBody](ctx)
	if validationErr != nil {
		return ctx.JSON(api.DetailedError(validationErr.Error, validationErr.Details))
	}

	slug := body.Slug
	if slug == "" {
		randomSlug, err := gonanoid.New(defaultSlugLength)
		if err != nil {
			slog.Error("Failed to generate random slug", "err", err)
			return ctx.JSON(api.Err(api.ErrInternalServer, "Something went wrong"))
		}
		slug = randomSlug
	}

	if slugBlacklist[slug] {
		return ctx.JSON(api.Err(api.ErrForbidden, "Slug is blacklisted"))
	}

	var app *config.AppConfig

	apiKey := ctx.Request().Header.Get("X-API-Key")
	if apiKey == "" {
		app = c.cfg.Public
	} else {
		app = c.cfg.AppByAPIKey[apiKey]
	}

	if app == nil || !app.Enabled {
		return ctx.JSON(api.Err(api.ErrUnauthorized, "Invalid API key"))
	}

	domain := ctx.Request().Host

	slog.Info("Creating link", "domain", domain, "slug", slug, "url", body.OriginalURL)

	ttlInSecs := *body.TTL
	var ttl time.Duration

	if ttlInSecs != 0 {
		ttl = time.Duration(ttlInSecs) * time.Second
	}

	if app.MaxDurationSec != 0 && ttlInSecs > app.MaxDurationSec {
		return ctx.JSON(api.Err(api.ErrBadRequest, fmt.Sprintf("TTL too high, max is %d", app.MaxDurationSec)))
	}

	if app.MinDurationSec != 0 && ttlInSecs < app.MinDurationSec {
		return ctx.JSON(api.Err(api.ErrBadRequest, fmt.Sprintf("TTL too low, min is %d", app.MinDurationSec)))
	}

	link := models.Link{
		Slug:        slug,
		Domain:      domain,
		OriginalURL: body.OriginalURL,
		TTL:         ttlInSecs,
		URL:         fmt.Sprintf("https://%s/%s", domain, slug),
	}

	key := fmt.Sprintf("link:%s/%s", domain, slug)
	cmd := c.vkey.B().Set().Key(key).Value(body.OriginalURL).Nx().Ex(ttl).Build()
	res := c.vkey.Do(context.Background(), cmd)

	if err := res.Error(); err != nil {
		if valkey.IsValkeyNil(err) {
			return ctx.JSON(api.Err(api.ErrConflict, "Link already exists"))
		}
		slog.Error("Failed to set key", "err", err)
		return ctx.JSON(api.Err(api.ErrInternalServer, "Something went wrong"))
	}

	return ctx.JSON(http.StatusCreated, link)
}

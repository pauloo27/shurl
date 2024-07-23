package link_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/pauloo27/shurl/internal/config"
	"github.com/pauloo27/shurl/internal/server/api/link"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestRedirect(t *testing.T) {
	rdb := mockRedis()

	rdb.FlushDB(context.Background())
	rdb.SetNX(context.Background(), "link:localhost/hello", "http://example.com", 30*time.Second)
	rdb.SetNX(context.Background(), "link:127.0.0.1/world", "http://example.com/world", 30*time.Second)

	t.Run("Valid domain and slug pair", func(t *testing.T) {
		cfg := &config.Config{}

		rec, err := callRedirectHandler(cfg, rdb, "localhost", "hello")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusTemporaryRedirect, rec.Code)
		assert.Equal(t, "http://example.com", rec.Header().Get("Location"))
	})

	t.Run("Mismatched domain and slug pair", func(t *testing.T) {
		cfg := &config.Config{}

		rec, err := callRedirectHandler(cfg, rdb, "localhost", "world")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t,
			`{"error":"NOT_FOUND","detail":{"message":"Link not found"}}`,
			strings.TrimSpace(rec.Body.String()),
		)
	})

	t.Run("Slug not found", func(t *testing.T) {
		cfg := &config.Config{}

		rec, err := callRedirectHandler(cfg, rdb, "localhost", "slug")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t,
			`{"error":"NOT_FOUND","detail":{"message":"Link not found"}}`,
			strings.TrimSpace(rec.Body.String()),
		)
	})
}

func TestRdbIsClosed(t *testing.T) {
	rdb := mockRedis()
	cfg := &config.Config{}

	err := rdb.Close()
	assert.NoError(t, err)

	t.Run("Rdb is closed", func(t *testing.T) {
		rec, err := callRedirectHandler(cfg, rdb, "localhost", "slug")

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(
			t, `{"error":"INTERNAL_SERVER_ERROR","detail":{"message":"Something went wrong"}}`,
			strings.TrimSpace(rec.Body.String()),
		)
	})
}

func callRedirectHandler(
	cfg *config.Config, rdb *redis.Client,
	domain, slug string,
) (*httptest.ResponseRecorder, error) {
	path := fmt.Sprintf("/%s", slug)

	e := echo.New()
	req := httptest.NewRequest("GET", path, nil)
	req.Host = domain
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetPath(path)
	ctx.SetParamNames("slug")
	ctx.SetParamValues(slug)
	c := link.NewLinkController(cfg, rdb)
	err := c.Redirect(ctx)
	return rec, err
}

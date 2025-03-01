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
	"github.com/stretchr/testify/assert"
	"github.com/valkey-io/valkey-go"
)

func TestRedirect(t *testing.T) {
	vkey := mockValkey()

	mustDo := func(cmd valkey.Completed) {
		res := vkey.Do(context.Background(), cmd)
		assert.NoError(t, res.Error())
	}

	flushCmd := vkey.B().Flushdb().Build()
	setHello := vkey.B().Set().Key("link:localhost/hello").Value("http://example.com").Nx().Ex(30 * time.Second).Build()
	setWorld := vkey.B().Set().Key("link:127.0.0.1/world").Value("http://example.com/world").Nx().Ex(30 * time.Second).Build()

	mustDo(flushCmd)
	mustDo(setHello)
	mustDo(setWorld)

	t.Run("Valid domain and slug pair", func(t *testing.T) {
		cfg := &config.Config{}

		rec, err := callRedirectHandler(cfg, vkey, "localhost", "hello")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusTemporaryRedirect, rec.Code)
		assert.Equal(t, "http://example.com", rec.Header().Get("Location"))
	})

	t.Run("Mismatched domain and slug pair", func(t *testing.T) {
		cfg := &config.Config{}

		rec, err := callRedirectHandler(cfg, vkey, "localhost", "world")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t,
			`{"error":"NOT_FOUND","detail":{"message":"Link not found"}}`,
			strings.TrimSpace(rec.Body.String()),
		)
	})

	t.Run("Slug not found", func(t *testing.T) {
		cfg := &config.Config{}

		rec, err := callRedirectHandler(cfg, vkey, "localhost", "slug")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t,
			`{"error":"NOT_FOUND","detail":{"message":"Link not found"}}`,
			strings.TrimSpace(rec.Body.String()),
		)
	})
}

func callRedirectHandler(
	cfg *config.Config, vkey valkey.Client,
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
	c := link.NewLinkController(cfg, vkey)
	err := c.Redirect(ctx)
	return rec, err
}

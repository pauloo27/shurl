package link_test

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/pauloo27/shurl/internal/config"
	"github.com/pauloo27/shurl/internal/mocker"
	"github.com/pauloo27/shurl/internal/server/api/link"
	"github.com/stretchr/testify/assert"
)

func callRedirectHandler(cfg *config.Config, domain, slug string) (*mocker.Response, error) {
	data := mocker.RequestData{
		Path: "/" + slug,
		Host: domain,
		URLParams: map[string]string{
			"slug": slug,
		},
		Method: "GET",
		Config: cfg,
		Rdb:    rdb,
	}
	return mocker.CallHandler(link.Redirect, data)
}

func TestRedirect(t *testing.T) {
	rdb.FlushDB(context.Background())
	rdb.SetNX(context.Background(), "link:localhost/hello", "http://example.com", 30*time.Second)
	rdb.SetNX(context.Background(), "link:127.0.0.1/world", "http://example.com/world", 30*time.Second)

	t.Run("Valid domain and slug pair", func(t *testing.T) {
		cfg := &config.Config{}

		res, err := callRedirectHandler(cfg, "localhost", "hello")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusTemporaryRedirect, res.StatusCode)
		assert.Equal(t, "http://example.com", res.Header.Get("Location"))
	})

	t.Run("Mismatched domain and slug pair", func(t *testing.T) {
		cfg := &config.Config{}

		res, err := callRedirectHandler(cfg, "localhost", "world")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, res.StatusCode)
		assert.Equal(t, `{"error":"NOT_FOUND","detail":{"message":"Link not found"}}`, strings.TrimSpace(res.StringBody))
	})

	t.Run("Slug not found", func(t *testing.T) {
		cfg := &config.Config{}

		res, err := callRedirectHandler(cfg, "localhost", "slug")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, res.StatusCode)
		assert.Equal(t, `{"error":"NOT_FOUND","detail":{"message":"Link not found"}}`, strings.TrimSpace(res.StringBody))
	})
}

func TestRdbIsClosed(t *testing.T) {
	err := rdb.Close()
	assert.NoError(t, err)

	t.Run("Rdb is closed", func(t *testing.T) {
		cfg := &config.Config{}

		res, err := callRedirectHandler(cfg, "localhost", "slug")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
		assert.Equal(
			t,
			`{"error":"INTERNAL_SERVER_ERROR","detail":{"message":"Something went wrong"}}`,
			strings.TrimSpace(res.StringBody),
		)
	})

	rdb = mocker.MakeRedictMock()
}

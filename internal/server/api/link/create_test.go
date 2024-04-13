package link_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/pauloo27/shurl/internal/config"
	"github.com/pauloo27/shurl/internal/mocker"
	"github.com/pauloo27/shurl/internal/server/api/link"
	"github.com/pauloo27/shurl/internal/server/validator"
	"github.com/stretchr/testify/assert"
)

func unmarshalAndValidate(raw string) (link.CreateLinkBody, bool) {
	r := httptest.NewRequest("POST", "/api/link", bytes.NewBufferString(raw))
	return validator.MustGetBody[link.CreateLinkBody](httptest.NewRecorder(), r)
}

func TestValidateBody(t *testing.T) {
	t.Run("Empty json", func(t *testing.T) {
		raw := `{}`
		_, ok := unmarshalAndValidate(raw)
		assert.False(t, ok)
	})

	t.Run("All fields present and valid", func(t *testing.T) {
		raw := `{"slug":"slug","domain":"domain","original_url":"http://example.com","ttl":1}`
		_, ok := unmarshalAndValidate(raw)
		assert.True(t, ok)
	})

	t.Run("Slug not present", func(t *testing.T) {
		raw := `{"domain":"domain","original_url":"http://example.com","ttl":1}`
		_, ok := unmarshalAndValidate(raw)
		assert.True(t, ok)
	})

	t.Run("Slug not present", func(t *testing.T) {
		raw := `{"domain":"domain","original_url":"http://example.com","ttl":1}`
		_, ok := unmarshalAndValidate(raw)
		assert.True(t, ok)
	})

	t.Run("Slug invalid", func(t *testing.T) {
		raw := `{"slug": "x", "domain":"domain","original_url":"http://example.com","ttl":1}`
		_, ok := unmarshalAndValidate(raw)
		assert.False(t, ok)
	})

	t.Run("Minimal required", func(t *testing.T) {
		raw := `{"original_url":"http://example.com","ttl":1}`
		_, ok := unmarshalAndValidate(raw)
		assert.True(t, ok)
	})

	t.Run("Original URL invalid", func(t *testing.T) {
		raw := `{"original_url":"example.com","ttl":1}`
		_, ok := unmarshalAndValidate(raw)
		assert.False(t, ok)
	})

	t.Run("TTL not present", func(t *testing.T) {
		raw := `{"original_url":"http://example.com"}`
		_, ok := unmarshalAndValidate(raw)
		assert.False(t, ok)
	})
}

func callCreateHandler(cfg *config.Config, apiKey, body string) (*mocker.Response, error) {
	headers := make(http.Header)
	headers.Set("X-API-Key", apiKey)

	data := mocker.RequestData{
		Body:    body,
		Headers: headers,
		Path:    "/api/link",
		Method:  "POST",
		Config:  mocker.MakeConfigMock(cfg),
		Rdb:     rdb,
	}
	return mocker.CallHandler(link.Create, data)
}

const (
	secretApiKey = "SECRET_API_KEY"
)

func TestAuthorization(t *testing.T) {
	t.Run("With no api key and public api disabled", func(t *testing.T) {
		config := &config.Config{
			Public: &config.AppConfig{
				Enabled: false,
			},
		}
		res, err := callCreateHandler(config, "", `{"original_url": "http://google.com", "ttl": 20}`)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, http.StatusUnauthorized, res.Status)
		assert.Equal(
			t,
			`{"detail":{"message":"Invalid API key"},"error":"UNAUTHORIZED"}`,
			strings.TrimSpace(res.Body),
		)
	})

	t.Run("With invalid api key", func(t *testing.T) {
		app := &config.AppConfig{
			APIKey:  secretApiKey,
			Enabled: true,
		}

		config := &config.Config{
			Public: &config.AppConfig{
				Enabled: false,
			},
			Apps: map[string]*config.AppConfig{
				"my-app": app,
			},
		}
		res, err := callCreateHandler(config, "wrong-api-key", `{"original_url": "http://google.com", "ttl": 20}`)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, http.StatusUnauthorized, res.Status)
		assert.Equal(
			t,
			`{"detail":{"message":"Invalid API key"},"error":"UNAUTHORIZED"}`,
			strings.TrimSpace(res.Body),
		)
	})

	t.Run("With valid api key, but disabled app", func(t *testing.T) {
		app := &config.AppConfig{
			APIKey:  secretApiKey,
			Enabled: false,
		}

		config := &config.Config{
			Public: &config.AppConfig{
				Enabled: false,
			},
			Apps: map[string]*config.AppConfig{
				"my-app": app,
			},
		}
		res, err := callCreateHandler(config, secretApiKey, `{"original_url": "http://google.com", "ttl": 20}`)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, http.StatusUnauthorized, res.Status)
		assert.Equal(
			t,
			`{"detail":{"message":"Invalid API key"},"error":"UNAUTHORIZED"}`,
			strings.TrimSpace(res.Body),
		)
	})

	t.Run("With public api enabled", func(t *testing.T) {
		config := &config.Config{
			Public: &config.AppConfig{
				Enabled:        true,
				AllowedDomains: []string{"localhost"},
			},
		}
		res, err := callCreateHandler(config, "", `{"original_url": "http://google.com", "ttl": 20}`)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, http.StatusCreated, res.Status)
	})

	t.Run("With valid api key", func(t *testing.T) {
		app := &config.AppConfig{
			APIKey:         secretApiKey,
			Enabled:        true,
			AllowedDomains: []string{"localhost"},
		}

		config := &config.Config{
			Public: &config.AppConfig{
				Enabled: false,
			},
			Apps: map[string]*config.AppConfig{
				"my-app": app,
			},
		}
		res, err := callCreateHandler(config, secretApiKey, `{"original_url": "http://google.com", "ttl": 20}`)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, http.StatusCreated, res.Status)
	})

	t.Run("With domain not allowed", func(t *testing.T) {
		app := &config.AppConfig{
			APIKey:         secretApiKey,
			Enabled:        true,
			AllowedDomains: []string{"localhost"},
		}

		config := &config.Config{
			Public: &config.AppConfig{
				Enabled: false,
			},
			Apps: map[string]*config.AppConfig{
				"my-app": app,
			},
		}
		res, err := callCreateHandler(config, secretApiKey, `{"domain": "google.com", "original_url": "http://google.com", "ttl": 20}`)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, http.StatusForbidden, res.Status)
		assert.Equal(
			t,
			`{"detail":{"message":"Domain not allowed for this app"},"error":"FORBIDDEN"}`,
			strings.TrimSpace(res.Body),
		)
	})
}

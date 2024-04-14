package link_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/pauloo27/shurl/internal/config"
	"github.com/pauloo27/shurl/internal/mocker"
	"github.com/pauloo27/shurl/internal/models"
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
		cfg := &config.Config{
			Public: &config.AppConfig{
				Enabled: false,
			},
		}
		res, err := callCreateHandler(cfg, "", `{"original_url": "http://google.com", "ttl": 20}`)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, http.StatusUnauthorized, res.Status)
		assert.Equal(
			t,
			`{"error":"UNAUTHORIZED","detail":{"message":"Invalid API key"}}`,
			strings.TrimSpace(res.Body),
		)
	})

	t.Run("With public api enabled, but no domain", func(t *testing.T) {
		cfg := &config.Config{
			Public: &config.AppConfig{
				Enabled: true,
			},
		}
		res, err := callCreateHandler(cfg, "", `{"original_url": "http://google.com", "ttl": 20}`)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, http.StatusForbidden, res.Status)
		assert.Equal(
			t,
			`{"error":"FORBIDDEN","detail":{"message":"No allowed domains for this app"}}`,
			strings.TrimSpace(res.Body),
		)
	})

	t.Run("With invalid api key", func(t *testing.T) {
		app := &config.AppConfig{
			APIKey:  secretApiKey,
			Enabled: true,
		}

		cfg := &config.Config{
			Public: &config.AppConfig{
				Enabled: false,
			},
			Apps: map[string]*config.AppConfig{
				"my-app": app,
			},
		}
		res, err := callCreateHandler(cfg, "wrong-api-key", `{"original_url": "http://google.com", "ttl": 20}`)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, http.StatusUnauthorized, res.Status)
		assert.Equal(
			t,
			`{"error":"UNAUTHORIZED","detail":{"message":"Invalid API key"}}`,
			strings.TrimSpace(res.Body),
		)
	})

	t.Run("With valid api key, but disabled app", func(t *testing.T) {
		app := &config.AppConfig{
			APIKey:  secretApiKey,
			Enabled: false,
		}

		cfg := &config.Config{
			Public: &config.AppConfig{
				Enabled: false,
			},
			Apps: map[string]*config.AppConfig{
				"my-app": app,
			},
		}
		res, err := callCreateHandler(cfg, secretApiKey, `{"original_url": "http://google.com", "ttl": 20}`)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, http.StatusUnauthorized, res.Status)
		assert.Equal(
			t,
			`{"error":"UNAUTHORIZED","detail":{"message":"Invalid API key"}}`,
			strings.TrimSpace(res.Body),
		)
	})

	t.Run("With public api enabled", func(t *testing.T) {
		cfg := &config.Config{
			Public: &config.AppConfig{
				Enabled:        true,
				AllowedDomains: []string{"localhost"},
			},
		}
		res, err := callCreateHandler(cfg, "", `{"original_url": "http://google.com", "ttl": 20}`)
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

		cfg := &config.Config{
			Public: &config.AppConfig{
				Enabled: false,
			},
			Apps: map[string]*config.AppConfig{
				"my-app": app,
			},
		}
		res, err := callCreateHandler(cfg, secretApiKey, `{"original_url": "http://google.com", "ttl": 20}`)
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

		cfg := &config.Config{
			Public: &config.AppConfig{
				Enabled: false,
			},
			Apps: map[string]*config.AppConfig{
				"my-app": app,
			},
		}
		res, err := callCreateHandler(cfg, secretApiKey, `{"domain": "google.com", "original_url": "http://google.com", "ttl": 20}`)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, http.StatusForbidden, res.Status)
		assert.Equal(
			t,
			`{"error":"FORBIDDEN","detail":{"message":"Domain not allowed for this app"}}`,
			strings.TrimSpace(res.Body),
		)
	})
}

func TestInvalidData(t *testing.T) {
	t.Run("With invalid json", func(t *testing.T) {
		cfg := &config.Config{}
		res, err := callCreateHandler(cfg, "", `{"`)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, http.StatusBadRequest, res.Status)
		assert.Equal(
			t,
			`{"error":"BAD_REQUEST","detail":{"message":"unexpected EOF"}}`,
			strings.TrimSpace(res.Body),
		)
	})

	t.Run("With missing ttl", func(t *testing.T) {
		cfg := &config.Config{}
		res, err := callCreateHandler(cfg, "", `{"original_url": "https://google.com"}`)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, http.StatusUnprocessableEntity, res.Status)
		assert.Equal(
			t,
			`{"error":"VALIDATION_ERROR","detail":[{"field":"ttl","error":"required"}]}`,
			strings.TrimSpace(res.Body),
		)
	})

	t.Run("Rdb is closed", func(t *testing.T) {
		err := rdb.Close()
		assert.NoError(t, err)

		cfg := &config.Config{
			Public: &config.AppConfig{
				Enabled:        true,
				AllowedDomains: []string{"localhost"},
			},
		}

		res, err := callCreateHandler(cfg, "", `{"original_url": "https://google.com", "ttl": 20}`)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, http.StatusInternalServerError, res.Status)
		assert.Equal(
			t,
			`{"error":"INTERNAL_SERVER_ERROR","detail":{"message":"Something went wrong"}}`,
			strings.TrimSpace(res.Body),
		)

		rdb = mocker.MakeRedictMock()
	})
}

func TestCreation(t *testing.T) {
	rdb.FlushDB(context.Background())

	t.Run("With specified domain and slug", func(t *testing.T) {
		cfg := &config.Config{
			Public: &config.AppConfig{
				Enabled:        true,
				AllowedDomains: []string{"localhost"},
			},
		}
		res, err := callCreateHandler(cfg, "", `{"domain": "localhost", "slug": "hello", "original_url": "http://google.com", "ttl": 23}`)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, http.StatusCreated, res.Status)
		assert.Equal(
			t,
			`{"slug":"hello","domain":"localhost","original_url":"http://google.com","ttl":23}`,
			strings.TrimSpace(res.Body),
		)

		rdbRes := rdb.Get(context.Background(), "link:localhost/hello")
		assert.Equal(t, "http://google.com", rdbRes.Val())
	})

	t.Run("With specified domain and random slug", func(t *testing.T) {
		cfg := &config.Config{
			Public: &config.AppConfig{
				Enabled:        true,
				AllowedDomains: []string{"localhost"},
			},
		}
		res, err := callCreateHandler(cfg, "", `{"domain": "localhost", "original_url": "http://google.com", "ttl": 23}`)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, http.StatusCreated, res.Status)

		var link models.Link
		err = json.Unmarshal([]byte(res.Body), &link)
		assert.NoError(t, err)

		assert.Equal(t, "localhost", link.Domain)
		assert.Equal(t, "http://google.com", link.OriginalURL)
		assert.Equal(t, 23, link.TTL)
		assert.NotEmpty(t, link.Slug)

		slug := link.Slug

		rdbRes := rdb.Get(context.Background(), "link:localhost/"+slug)
		assert.Equal(t, "http://google.com", rdbRes.Val())
	})

	t.Run("With random slug and no domain", func(t *testing.T) {
		cfg := &config.Config{
			Public: &config.AppConfig{
				Enabled:        true,
				AllowedDomains: []string{"localhost"},
			},
		}
		res, err := callCreateHandler(cfg, "", `{"original_url": "http://google.com", "ttl": 23}`)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, http.StatusCreated, res.Status)

		var link models.Link
		err = json.Unmarshal([]byte(res.Body), &link)
		assert.NoError(t, err)

		assert.Equal(t, "localhost", link.Domain)
		assert.Equal(t, "http://google.com", link.OriginalURL)
		assert.Equal(t, 23, link.TTL)
		assert.NotEmpty(t, link.Slug)

		slug := link.Slug

		rdbRes := rdb.Get(context.Background(), "link:localhost/"+slug)
		assert.Equal(t, "http://google.com", rdbRes.Val())
	})

	t.Run("With slug and domain pair already in use", func(t *testing.T) {
		cfg := &config.Config{
			Public: &config.AppConfig{
				Enabled:        true,
				AllowedDomains: []string{"localhost"},
			},
		}
		res, err := callCreateHandler(cfg, "", `{"domain":"localhost", "slug": "flamengo", "original_url": "http://google.com", "ttl": 23}`)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, http.StatusCreated, res.Status)

		res, err = callCreateHandler(cfg, "", `{"domain":"localhost", "slug": "flamengo", "original_url": "http://bing.com", "ttl": 23}`)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, http.StatusConflict, res.Status)
		assert.Equal(
			t,
			`{"error":"CONFLICT","detail":{"message":"Link already exists"}}`,
			strings.TrimSpace(res.Body),
		)

		rdbRes := rdb.Get(context.Background(), "link:localhost/flamengo")
		assert.Equal(t, "http://google.com", rdbRes.Val())
	})
}

func TestBlacklistedSlugs(t *testing.T) {
	test := func(slug string) {
		cfg := &config.Config{
			Public: &config.AppConfig{
				Enabled: true,
			},
		}
		res, err := callCreateHandler(
			cfg,
			"",
			`{"slug": "`+slug+`", "original_url": "http://google.com", "ttl": 23}`,
		)
		assert.NoError(t, err)
		assert.NotNil(t, res)

		assert.Equal(t, http.StatusForbidden, res.Status)
		assert.Equal(
			t,
			`{"error":"FORBIDDEN","detail":{"message":"Slug is blacklisted"}}`,
			strings.TrimSpace(res.Body),
		)
	}

	for slug := range link.SlugBlacklist {
		t.Run("With blacklisted slug "+slug, func(t *testing.T) {
			test(slug)
		})
	}
}

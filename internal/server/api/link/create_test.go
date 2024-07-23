package link_test

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/pauloo27/shurl/internal/server/api/link"
	"github.com/pauloo27/shurl/internal/server/core/validator"
	"github.com/stretchr/testify/assert"
)

func unmarshalAndValidate(raw string) (link.CreateLinkBody, bool) {
	e := echo.New()
	req := httptest.NewRequest("POST", "/links", strings.NewReader(raw))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	body, err := validator.MustBindAndValidate[link.CreateLinkBody](ctx)
	return body, err == nil
}

func TestValidateBody(t *testing.T) {
	t.Run("Empty json", func(t *testing.T) {
		raw := `{}`
		_, ok := unmarshalAndValidate(raw)
		assert.False(t, ok)
	})

	t.Run("All fields present and valid", func(t *testing.T) {
		raw := `{"slug":"slug","original_url":"http://example.com","ttl":1}`
		_, ok := unmarshalAndValidate(raw)
		assert.True(t, ok)
	})

	t.Run("Slug not present", func(t *testing.T) {
		raw := `{"original_url":"http://example.com","ttl":1}`
		_, ok := unmarshalAndValidate(raw)
		assert.True(t, ok)
	})

	t.Run("Slug not present", func(t *testing.T) {
		raw := `{"original_url":"http://example.com","ttl":1}`
		_, ok := unmarshalAndValidate(raw)
		assert.True(t, ok)
	})

	t.Run("Slug invalid", func(t *testing.T) {
		raw := `{"slug": "x","original_url":"http://example.com","ttl":1}`
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

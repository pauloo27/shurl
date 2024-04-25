package router

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pauloo27/shurl/internal/server/api"
)

func TestWrap(t *testing.T) {
	handler := func(r *http.Request) api.Response {
		return api.Response{
			StatusCode: http.StatusTeapot,
			Header:     map[string][]string{"Content-Type": {"application/json"}},
			Body:       map[string]any{"message": "success"},
		}
	}

	wrappedHandler := wrap(handler)

	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusTeapot, rr.Code)

	expected := `{"message":"success"}`
	assert.Equal(t, expected, strings.TrimSpace(rr.Body.String()))
}

func TestWrapWithMissingResponse(t *testing.T) {
	handler := func(r *http.Request) api.Response {
		return api.Response{}
	}

	wrappedHandler := wrap(handler)

	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	expected := `{"error":"Missing response","detail":null}`
	assert.Equal(t, expected, strings.TrimSpace(rr.Body.String()))
}

func TestWrapWithNilBody(t *testing.T) {
	handler := func(r *http.Request) api.Response {
		return api.Response{
			StatusCode: http.StatusOK,
			Header:     map[string][]string{"Content-Type": {"application/json"}},
			Body:       nil,
		}
	}

	wrappedHandler := wrap(handler)

	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	assert.Equal(t, "", strings.TrimSpace(rr.Body.String()))
}

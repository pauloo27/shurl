package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoggerMiddleware(t *testing.T) {
	buf := strings.Builder{}

	h := slog.NewJSONHandler(&buf, nil)
	slog.SetDefault(slog.New(h))

	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "127.0.0.1"

	ctx := context.WithValue(context.Background(), middleware.RequestIDKey, "123")
	req = req.WithContext(ctx)

	res := httptest.NewRecorder()

	handler := loggerMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(res, req)

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")

	expectedLines := []map[string]any{
		{
			"level":       "INFO",
			"msg":         "Http request",
			"id":          "123",
			"remote_addr": "127.0.0.1",
			"method":      "GET",
		},
		{
			"level":  "INFO",
			"msg":    "Http response",
			"id":     "123",
			"status": float64(200),
		},
	}

	for i, line := range lines {
		var record map[string]any
		err := json.Unmarshal([]byte(line), &record)
		require.NoError(t, err)
		fmt.Println(line)

		expectedLine := expectedLines[i]

		for k, v := range expectedLine {
			assert.Equal(t, v, record[k])
		}
	}
}

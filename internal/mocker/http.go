package mocker

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/pauloo27/shurl/internal/config"
	"github.com/pauloo27/shurl/internal/ctx"
	"github.com/redis/go-redis/v9"
)

type Response struct {
	Status int
	Body   string
}

type RequestData struct {
	Path    string
	Method  string
	Body    string
	Headers http.Header

	// not request, per se, but needed for the handler
	Config *config.Config
	Rdb    *redis.Client
}

func CallHandler(handler http.HandlerFunc, data RequestData) (*Response, error) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(data.Method, data.Path, strings.NewReader(data.Body))
	r.Header = data.Headers

	cfg := MakeConfigMock(data.Config)
	rdb := data.Rdb

	services := &ctx.Services{
		Config: cfg,
		Rdb:    rdb,
	}

	r = r.WithContext(context.WithValue(r.Context(), ctx.ServicesKey, services))

	handler(w, r)

	rawBody, err := io.ReadAll(w.Result().Body)
	if err != nil {
		return nil, err
	}

	res := Response{
		Status: w.Result().StatusCode,
		Body:   string(rawBody),
	}

	return &res, nil
}

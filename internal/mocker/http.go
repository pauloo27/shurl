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
	Status  int
	Body    string
	Headers http.Header
}

type RequestData struct {
	Host      string
	Path      string
	Method    string
	Body      string
	Headers   http.Header
	URLParams map[string]string

	// not request, per se, but needed for the handler
	Config *config.Config
	Rdb    *redis.Client
}

func CallHandler(handler http.HandlerFunc, data RequestData) (*Response, error) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(data.Method, data.Path, strings.NewReader(data.Body))
	r.Header = data.Headers
	r.Host = data.Host

	services := &ctx.Services{
		Config: data.Config,
		Rdb:    data.Rdb,
	}

	r = r.WithContext(context.WithValue(r.Context(), ctx.ServicesKey, services))

	for k, v := range data.URLParams {
		r.SetPathValue(k, v)
	}

	handler(w, r)

	rawBody, err := io.ReadAll(w.Result().Body)
	if err != nil {
		return nil, err
	}

	res := Response{
		Status:  w.Result().StatusCode,
		Body:    string(rawBody),
		Headers: w.Result().Header,
	}

	return &res, nil
}

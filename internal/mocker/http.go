package mocker

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/pauloo27/shurl/internal/config"
	"github.com/pauloo27/shurl/internal/ctx"
	"github.com/pauloo27/shurl/internal/server/api"
	"github.com/pauloo27/shurl/internal/server/router"
	"github.com/redis/go-redis/v9"
)

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

type Response struct {
	*api.Response
	StringBody string
}

func CallHandler(handler router.WrappedHandler, data RequestData) (*Response, error) {
	r := httptest.NewRequest(data.Method, data.Path, strings.NewReader(data.Body))
	r.Header = data.Headers
	r.Host = data.Host

	providers := &ctx.Providers{
		Config: data.Config,
		Rdb:    data.Rdb,
		Logger: slog.Default(),
	}

	r = r.WithContext(context.WithValue(r.Context(), ctx.ProvidersKey, providers))

	for k, v := range data.URLParams {
		r.SetPathValue(k, v)
	}

	res := handler(r)

	encodedBody, err := json.Marshal(res.Body)
	if err != nil {
		return nil, err
	}

	return &Response{
		Response:   &res,
		StringBody: string(encodedBody),
	}, nil
}

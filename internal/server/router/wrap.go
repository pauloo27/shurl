package router

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/pauloo27/shurl/internal/server/api"
)

type WrappedHandler func(r *http.Request) api.Response

func wrap(handler WrappedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res := handler(r)

		// TODO: avoid clone?
		for k, v := range res.Header {
			w.Header().Set(k, v[0])
		}

		if res.StatusCode == 0 {
			slog.Error("Missing response status code", "path", r.URL.Path)
			res.StatusCode = http.StatusInternalServerError
		}

		w.WriteHeader(res.StatusCode)

		if res.Body != nil {
			_ = json.NewEncoder(w).Encode(res.Body)
		}
	}
}

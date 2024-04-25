package router

import (
	"encoding/json"
	"net/http"

	"github.com/pauloo27/shurl/internal/server/api"
)

type WrappedHandler func(r *http.Request) api.Response

func wrap(handler WrappedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res := handler(r)

		if res.StatusCode == 0 {
			res.StatusCode = http.StatusInternalServerError
			res.Body = api.Error[any]{
				Error:  "Missing response",
				Detail: nil,
			}
		}

		// TODO: avoid clone?
		for k, v := range res.Header {
			w.Header().Set(k, v[0])
		}

		w.WriteHeader(res.StatusCode)

		if res.Body != nil {
			_ = json.NewEncoder(w).Encode(res.Body)
		}
	}
}

package validator

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/lmittmann/tint"
	"github.com/pauloo27/shurl/internal/server/api"
)

func MustGetBody[T any](w http.ResponseWriter, r *http.Request) (T, bool) {
	var body T

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		slog.Error("Failed to decode body", tint.Err(err))
		if errors.Is(err, io.EOF) {
			api.Err(w, http.StatusBadRequest, api.BadRequestErr, "Missing body")
		} else {
			api.Err(w, http.StatusBadRequest, api.BadRequestErr, err.Error())
		}
		return body, false
	}

	validationErrors := Validate[T](body)

	if validationErrors == nil {
		return body, true
	}

	api.DetailedError(w, http.StatusUnprocessableEntity, api.ValidationErr, validationErrors)

	return body, false
}

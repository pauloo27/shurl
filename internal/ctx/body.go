package ctx

import (
	"encoding/json"
	"net/http"

	"github.com/pauloo27/shurl/internal/server/validator"
)

func MustGetBody[T any](w http.ResponseWriter, r *http.Request) (T, bool) {
	var body T

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return body, false
	}

	validationErrors := validator.Validate[T](body)

	if validationErrors == nil {
		return body, true
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnprocessableEntity)
	json.NewEncoder(w).Encode(validationErrors)

	return body, false
}

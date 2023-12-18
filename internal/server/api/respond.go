package api

import (
	"encoding/json"
	"net/http"
)

type ErrorType string

const (
	NotFoundErr       ErrorType = "NOT_FOUND"
	InternalServerErr ErrorType = "INTERNAL_SERVER_ERROR"
	BadRequestErr     ErrorType = "BAD_REQUEST"
	ValidationErr     ErrorType = "VALIDATION_ERROR"
)

func Created[T any](w http.ResponseWriter, detail T) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(detail)
}

func Err(
	w http.ResponseWriter,
	statusCode int,
	error ErrorType,
	message string,
) {
	DetailedError(w, statusCode, error, map[string]string{
		"message": message,
	})
}

func DetailedError[T any](
	w http.ResponseWriter,
	statusCode int,
	error ErrorType,
	detail T,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"error":  error,
		"detail": detail,
	})
}

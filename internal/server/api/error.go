package api

import "net/http"

type ErrorType struct {
	Name       string `json:"name"`
	StatusCode int    `json:"status_code"`
}

var (
	NotFoundErr       = ErrorType{"NOT_FOUND", http.StatusNotFound}
	InternalServerErr = ErrorType{"INTERNAL_SERVER_ERROR", http.StatusInternalServerError}
	ForbiddenErr      = ErrorType{"FORBIDDEN", http.StatusForbidden}
	ConflictErr       = ErrorType{"CONFLICT", http.StatusConflict}
	BadRequestErr     = ErrorType{"BAD_REQUEST", http.StatusBadRequest}
	ValidationErr     = ErrorType{"VALIDATION_ERROR", http.StatusUnprocessableEntity}
	UnauthorizedErr   = ErrorType{"UNAUTHORIZED", http.StatusUnauthorized}
)

type Error[T any] struct {
	Error  string `json:"error"`
	Detail any    `json:"detail"`
}

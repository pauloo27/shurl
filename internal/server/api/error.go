package api

import "net/http"

type ErrorType struct {
	Name       string `json:"name"`
	StatusCode int    `json:"status_code"`
}

var (
	ErrNotFound       = ErrorType{"NOT_FOUND", http.StatusNotFound}
	ErrInternalServer = ErrorType{"INTERNAL_SERVER_ERROR", http.StatusInternalServerError}
	ErrForbidden      = ErrorType{"FORBIDDEN", http.StatusForbidden}
	ErrConflict       = ErrorType{"CONFLICT", http.StatusConflict}
	ErrBadRequest     = ErrorType{"BAD_REQUEST", http.StatusBadRequest}
	ErrValidation     = ErrorType{"VALIDATION_ERROR", http.StatusUnprocessableEntity}
	ErrUnauthorized   = ErrorType{"UNAUTHORIZED", http.StatusUnauthorized}
	ErrNotImplemented = ErrorType{"NOT_IMPLEMENTED", http.StatusNotImplemented}
)

type Error[T any] struct {
	Error  string `json:"error" example:"NOT_FOUND"`
	Detail T      `json:"detail"`
}

func Err(
	error ErrorType,
	message string,
) (status int, body Error[map[string]string]) {
	return DetailedError(error, map[string]string{
		"message": message,
	})
}

func DetailedError[T any](
	error ErrorType,
	detail T,
) (status int, body Error[T]) {
	header := make(http.Header)
	header.Set("Content-Type", "application/json")

	return error.StatusCode, Error[T]{
		Error:  error.Name,
		Detail: detail,
	}
}

// types used for better swagger docs
type BadRequestError struct {
	Error  string            `json:"error" example:"BAD_REQUEST"`
	Detail map[string]string `json:"detail" example:"message:Error message"`
}

type ValidationErrorDetail struct {
	Field string `json:"field" example:"username"`
	Error string `json:"error" example:"required"`
}

type ValidationError struct {
	Error  string                  `json:"error" example:"VALIDATION_ERROR"`
	Detail []ValidationErrorDetail `json:"detail"`
}

type InternalServerError struct {
	Error  string            `json:"error" example:"INTERNAL_SERVER_ERROR"`
	Detail map[string]string `json:"detail" example:"message:Error message"`
}

type NotImplementedError struct {
	Error  string            `json:"error" example:"NOT_IMPLEMENTED"`
	Detail map[string]string `json:"detail" example:"message:Error message"`
}

type ConflictError struct {
	Error  string            `json:"error" example:"CONFLICT"`
	Detail map[string]string `json:"detail" example:"message:Error message"`
}

type UnauthorizedError struct {
	Error  string            `json:"error" example:"UNAUTHORIZED"`
	Detail map[string]string `json:"detail" example:"message:Error message"`
}

type ForbiddenError struct {
	Error  string            `json:"error" example:"FORBIDDEN"`
	Detail map[string]string `json:"detail" example:"message:Error message"`
}

type NotFoundError struct {
	Error  string            `json:"error" example:"NOT_FOUND"`
	Detail map[string]string `json:"detail" example:"message:Error message"`
}

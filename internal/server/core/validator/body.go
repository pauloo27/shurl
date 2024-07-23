package validator

import (
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/pauloo27/shurl/internal/server/api"
)

type APIBodyValidationError struct {
	Details any
	Error   api.ErrorType
}

func MustBindAndValidate[T any](ctx echo.Context) (T, *APIBodyValidationError) {
	var payload T

	err := ctx.Bind(&payload)
	if err != nil {
		slog.Error("Failed to decode payload", "err", err)
		return payload, &APIBodyValidationError{
			Error:   api.ErrBadRequest,
			Details: map[string]string{"message": "Invalid payload"},
		}
	}

	validationErrors := Validate(payload)

	if len(validationErrors) > 0 {
		return payload, &APIBodyValidationError{
			Error:   api.ErrValidation,
			Details: validationErrors,
		}
	}

	return payload, nil
}

package health

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

type HealthStatus struct {
	Valkey bool `json:"valkey"`
}

// Health godoc
//
//	@Summary		Get health status
//	@Description	Get the health status of the server
//	@Tags			health
//	@Produce		json
//	@Success		200	{object}	HealthStatus
//	@Failure		500	{object}	HealthStatus
//	@Router			/healthz [get]
func (c *HealthController) Health(ctx echo.Context) error {
	ok := true
	status := HealthStatus{
		Valkey: true,
	}

	cmd := c.vkey.B().Ping().Build()
	if res := c.vkey.Do(context.Background(), cmd); res.Error() != nil {
		status.Valkey = false
		ok = false
	}

	if !ok {
		slog.Error("Health check failed", "status", status)
		return ctx.JSON(http.StatusInternalServerError, status)
	}

	return ctx.JSON(http.StatusOK, status)
}

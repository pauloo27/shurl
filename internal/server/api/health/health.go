package health

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

type HealthStatus struct {
	Rdb bool `json:"rdb"`
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
		Rdb: true,
	}

	rdbRes := c.rdb.Ping(context.Background())
	if rdbRes.Err() != nil {
		status.Rdb = false
		ok = false
	}

	if !ok {
		slog.Error("Health check failed", "status", status)
		return ctx.JSON(http.StatusInternalServerError, status)
	}

	return ctx.JSON(http.StatusOK, status)
}

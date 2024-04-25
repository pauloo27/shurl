package health

import (
	"context"
	"net/http"

	"github.com/pauloo27/shurl/internal/ctx"
	"github.com/pauloo27/shurl/internal/server/api"
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
func Health(r *http.Request) api.Response {
	providers := ctx.GetProviders(r.Context())
	log := providers.Logger

	ok := true
	status := HealthStatus{
		Rdb: true,
	}

	rdbRes := providers.Rdb.Ping(context.Background())
	if rdbRes.Err() != nil {
		status.Rdb = false
		ok = false
	}

	if !ok {
		log.Error("Health check failed", "status", status)
		return api.DetailedError(api.InternalServerErr, status)
	}

	return api.Ok(status)
}

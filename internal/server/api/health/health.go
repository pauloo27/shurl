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
func Health(w http.ResponseWriter, r *http.Request) {
	providers := ctx.GetProviders(r.Context())

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
		api.DetailedError(w, api.InternalServerErr, status)
		return
	}

	api.Ok(w, status)
}

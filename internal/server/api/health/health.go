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

func Health(w http.ResponseWriter, r *http.Request) {
	services := ctx.GetServices(r.Context())

	ok := true
	status := HealthStatus{
		Rdb: true,
	}

	rdbRes := services.Rdb.Ping(context.Background())
	if rdbRes.Err() != nil {
		status.Rdb = false
		ok = false
	}

	if !ok {
		api.DetailedError(w, http.StatusInternalServerError, api.InternalServerErr, status)
		return
	}

	api.Created(w, status)
}

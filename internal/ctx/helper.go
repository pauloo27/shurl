package ctx

import (
	"context"
)

func GetServices(ctx context.Context) *Services {
	return ctx.Value(ServicesKey).(*Services)
}

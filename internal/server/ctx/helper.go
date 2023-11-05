package ctx

import (
	"context"

	"github.com/pauloo27/shurl/internal/app"
)

func GetApp(ctx context.Context) *app.App {
	return ctx.Value(AppKey).(*app.App)
}

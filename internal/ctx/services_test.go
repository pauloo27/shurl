package ctx_test

import (
	"context"
	"testing"

	"github.com/pauloo27/shurl/internal/ctx"
	"github.com/stretchr/testify/assert"
)

func TestGetServicesFromContext(t *testing.T) {
	services := &ctx.Services{}

	c := context.Background()
	c = context.WithValue(c, ctx.ServicesKey, services)

	anotherService := ctx.GetServices(c)

	assert.Equal(t, services, anotherService)
}

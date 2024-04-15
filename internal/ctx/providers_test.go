package ctx_test

import (
	"context"
	"testing"

	"github.com/pauloo27/shurl/internal/ctx"
	"github.com/stretchr/testify/assert"
)

func TestGetProvidersFromContext(t *testing.T) {
	providers := &ctx.Providers{}

	c := context.Background()
	c = context.WithValue(c, ctx.ProvidersKey, providers)

	anotherProviders := ctx.GetProviders(c)

	assert.Equal(t, providers, anotherProviders)
}

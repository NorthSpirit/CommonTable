package externalapi

import (
	"context"

	"CommonTable/internal/repository"
)

type contextKey string

const apiKeyContextKey contextKey = "externalApiKey"

func withApiKey(ctx context.Context, key repository.RestaurantApiKey) context.Context {
	return context.WithValue(ctx, apiKeyContextKey, key)
}

func apiKeyFromContext(ctx context.Context) repository.RestaurantApiKey {
	key, _ := ctx.Value(apiKeyContextKey).(repository.RestaurantApiKey)
	return key
}

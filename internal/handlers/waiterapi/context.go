package waiterapi

import (
	"context"

	"CommonTable/internal/services/authservices"
)

type contextKey string

const claimsContextKey contextKey = "waiterClaims"

func withClaims(ctx context.Context, claims *authservices.EmployeeClaims) context.Context {
	return context.WithValue(ctx, claimsContextKey, claims)
}

func claimsFromContext(ctx context.Context) *authservices.EmployeeClaims {
	claims, _ := ctx.Value(claimsContextKey).(*authservices.EmployeeClaims)
	return claims
}

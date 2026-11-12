package externalapi

import (
	"errors"
	"net/http"
	"strings"

	"CommonTable/internal/services/externalorderservices"
)

func RequireApiKey(keys *externalorderservices.ApiKeyService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				writeError(w, http.StatusUnauthorized, "missing or invalid authorization header")
				return
			}

			key, err := keys.Authenticate(r.Context(), strings.TrimPrefix(authHeader, "Bearer "))
			if err != nil {
				if errors.Is(err, externalorderservices.ErrInvalidApiKey) {
					writeError(w, http.StatusUnauthorized, "invalid or revoked api key")
					return
				}
				writeError(w, http.StatusInternalServerError, "failed to check api key")
				return
			}

			next.ServeHTTP(w, r.WithContext(withApiKey(r.Context(), key)))
		})
	}
}

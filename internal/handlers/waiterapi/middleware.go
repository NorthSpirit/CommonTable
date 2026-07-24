package waiterapi

import (
	"net/http"
	"strings"
	"time"

	"CommonTable/internal/services/authservices"
)

// tokenTTL is intentionally short — the waiter app's tokens use a sliding
// expiry (see RequireWaiterOrManager below), so the window only needs to
// cover the gap between requests, not the whole shift. A lost or stolen
// device stops refreshing and the session lapses within this window.
const tokenTTL = time.Hour

// RequireWaiterOrManager validates a bearer token from the Authorization
// header, requires the employee's role to be "waiter" or "manager", and
// — on success — reissues a fresh token with a renewed expiry, returned
// via the X-Refreshed-Token response header. The client should store
// whichever token it most recently received and use that on the next
// call, so an actively-used session never expires mid-shift, while an
// idle/lost device's session lapses within one tokenTTL window.
func RequireWaiterOrManager(authService *authservices.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				writeError(w, http.StatusUnauthorized, "missing or invalid authorization header")
				return
			}
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			claims, err := authService.ValidateEmployeeToken(tokenString)
			if err != nil {
				writeError(w, http.StatusUnauthorized, "invalid or expired token")
				return
			}

			if claims.Role != "waiter" && claims.Role != "manager" {
				writeError(w, http.StatusForbidden, "this endpoint is for waiters and managers only")
				return
			}

			if refreshed, err := authService.RefreshEmployeeToken(claims, tokenTTL); err == nil {
				w.Header().Set("X-Refreshed-Token", refreshed)
			}

			next.ServeHTTP(w, r.WithContext(withClaims(r.Context(), claims)))
		})
	}
}

package waiterapi

import (
	"net/http"
	"strings"
	"time"

	"CommonTable/internal/services/authservices"
)

// Waiters token, in nature differs from those of the manager or admin,
// At first, I was trying to recall, what it was like for me to work as a waiter and calculate as safe time
// for the token to expire, but there is not really a good time - if a shift starts at 06:00
// a short fixed expiry would run out in the morning, around the time first customers arrive
// A slightly longer one, maybe 6-7ish hours could run out at the worst possible time - during the lunch peak hours
// which would have been horrible, on the other hand, I did think, it's an absolutely must for security reasons
// Eventually got the idea that, waiter's token shouldn't expire while waiter is actually working
// But expire relatively soon, if device is not being used, in the scenarios, where someone forgets to log off.
// So in the end, settled to 1 hour, with it refreshing easily
const tokenTTL = time.Hour

// RequireWaiterOrManager Method validates the bearer token from the Authorization header and requires the specific role from the user
// it must either be waiter or a manager and it is expected that the client will store the new token themselves
func RequireWaiterOrManager(authService *authservices.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// First it parses out the authorization token from the header part.
			authHeader := r.Header.Get("Authorization")

			// The header must start with "Bearer " (that is how the token is sent), otherwise returns 401
			if !strings.HasPrefix(authHeader, "Bearer ") {
				writeError(w, http.StatusUnauthorized, "missing or invalid authorization header")
				return
			}

			// Removes the bearer prefix to check the actual token
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			// Second, checks the token's signature (nobody messed with it) and expiration date, both come from the token itself, not from the DB
			claims, err := authService.ValidateEmployeeToken(tokenString)
			if err != nil {
				writeError(w, http.StatusUnauthorized, "invalid or expired token")
				return
			}

			// Third, check the role of the token owner.
			// Basically, the app is meant to be used for waiters, since orders is not the concern of the chef
			// On other hand, it is possible that in a restaurant, manager has to do basically everything
			// Waiting not excluded and it would be irrational for them to have another account just for the app
			if claims.Role != "waiter" && claims.Role != "manager" {
				writeError(w, http.StatusForbidden, "this endpoint is for waiters and managers only")
				return
			}

			// Eventually refreshes the token, by generating another, hour-lasting token to be reused again
			// Refreshing generates the tokens without having to look up into the DB
			// Since person's ID, restaurant ID, person's role and name is gotten from the token itself
			// Served with claims
			// If refreshing fails for some reason, the request still goes on, the client just keeps the old token
			if refreshed, err := authService.RefreshEmployeeToken(claims, tokenTTL); err == nil {
				w.Header().Set("X-Refreshed-Token", refreshed)
			}

			// Everything passed, so the request goes on to the next handler in the chain (like Create)
			// The claims travel with it: withClaims puts them into a copy of the request's context,
			// and the handler reads them back with claimsFromContext
			// If this wasn't called, the request would end here, which is what the returns above do
			next.ServeHTTP(w, r.WithContext(withClaims(r.Context(), claims)))
		})
	}
}

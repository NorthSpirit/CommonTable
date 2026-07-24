package kitchen

import (
	"net/http"

	"CommonTable/internal/services/authservices"
)

// RequireKitchenStaff protects a route group behind a valid employee
// session belonging to a chef or manager — a waiter's credentials are
// rejected here, same enforcement pattern as chef.RequireChef.
func RequireKitchenStaff(authService *authservices.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(sessionCookieName)
			if err != nil {
				http.Redirect(w, r, "/kitchen/login", http.StatusFound)
				return
			}

			claims, err := authService.ValidateEmployeeToken(cookie.Value)
			if err != nil {
				clearSessionCookie(w)
				http.Redirect(w, r, "/kitchen/login", http.StatusFound)
				return
			}

			if claims.Role != "chef" && claims.Role != "manager" {
				clearSessionCookie(w)
				http.Redirect(w, r, "/kitchen/login", http.StatusFound)
				return
			}

			next.ServeHTTP(w, r.WithContext(withClaims(r.Context(), claims)))
		})
	}
}

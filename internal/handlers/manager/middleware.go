package manager

import (
	"net/http"

	"CommonTable/internal/services/authservices"
)

// RequireManager protects a route group behind a valid employee session
// cookie belonging to an employee with role "manager" — a chef's or
// waiter's otherwise-valid token is rejected here, same enforcement
// pattern as chef.RequireChef and kitchen.RequireKitchenStaff.
func RequireManager(authService *authservices.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(sessionCookieName)
			if err != nil {
				http.Redirect(w, r, "/manager/login", http.StatusFound)
				return
			}

			claims, err := authService.ValidateEmployeeToken(cookie.Value)
			if err != nil {
				clearSessionCookie(w)
				http.Redirect(w, r, "/manager/login", http.StatusFound)
				return
			}

			if claims.Role != "manager" {
				clearSessionCookie(w)
				http.Redirect(w, r, "/manager/login", http.StatusFound)
				return
			}

			next.ServeHTTP(w, r.WithContext(withClaims(r.Context(), claims)))
		})
	}
}

package staff

import (
	"net/http"

	"CommonTable/internal/services/authservices"
)

// RequireStaff protects a route group behind a valid employee session
// cookie belonging to an employee with role "waiter" — a chef's or
// manager's otherwise-valid token is rejected here, same enforcement
// pattern as chef.RequireChef and manager.RequireManager.
func RequireStaff(authService *authservices.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(sessionCookieName)
			if err != nil {
				http.Redirect(w, r, "/staff/login", http.StatusFound)
				return
			}

			claims, err := authService.ValidateEmployeeToken(cookie.Value)
			if err != nil {
				clearSessionCookie(w)
				http.Redirect(w, r, "/staff/login", http.StatusFound)
				return
			}

			if claims.Role != "waiter" {
				clearSessionCookie(w)
				http.Redirect(w, r, "/staff/login", http.StatusFound)
				return
			}

			next.ServeHTTP(w, r.WithContext(withClaims(r.Context(), claims)))
		})
	}
}

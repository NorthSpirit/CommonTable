package chef

import (
	"net/http"

	"CommonTable/internal/services/authservices"
)

// RequireChef protects a route group behind a valid employee session cookie
// belonging to an employee with role "chef". Any other role — including a
// perfectly valid waiter/manager token — is rejected here, which is the
// enforcement point for "a waiter can't use the chef web page."
func RequireChef(authService *authservices.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(sessionCookieName)
			if err != nil {
				http.Redirect(w, r, "/chef/login", http.StatusFound)
				return
			}

			claims, err := authService.ValidateEmployeeToken(cookie.Value)
			if err != nil {
				clearSessionCookie(w)
				http.Redirect(w, r, "/chef/login", http.StatusFound)
				return
			}

			if claims.Role != "chef" {
				clearSessionCookie(w)
				http.Redirect(w, r, "/chef/login", http.StatusFound)
				return
			}

			next.ServeHTTP(w, r.WithContext(withClaims(r.Context(), claims)))
		})
	}
}

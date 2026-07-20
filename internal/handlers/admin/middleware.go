package admin

import (
	"net/http"

	"CommonTable/internal/services/adminservices"
)

// RequireAdmin protects a route group behind a valid admin session cookie.
// On failure (missing or invalid cookie), redirects to the login page
// rather than returning a bare 401 — this is a browser-navigated panel,
// not a JSON API, so a redirect is the correct UX.
func RequireAdmin(adminService *adminservices.AdminService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(sessionCookieName)
			if err != nil {
				http.Redirect(w, r, "/admin/login", http.StatusFound)
				return
			}

			_, err = adminService.ValidateAdminToken(cookie.Value)
			if err != nil {
				clearSessionCookie(w)
				http.Redirect(w, r, "/admin/login", http.StatusFound)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

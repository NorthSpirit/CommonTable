package terminal

import "net/http"

const sessionCookieName = "terminal_session"

func setSessionCookie(w http.ResponseWriter, restaurantID string) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    restaurantID,
		Path:     "/terminal",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   60 * 60 * 12, // A terminal stays "logged in" for a long shift.
	})
}

func clearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/terminal",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})
}

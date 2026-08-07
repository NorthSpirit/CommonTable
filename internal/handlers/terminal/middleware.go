package terminal

import (
	"context"
	"net/http"
	"strconv"
)

type contextKey string

const restaurantIDKey contextKey = "terminalRestaurantID"

func RequireTerminalSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(sessionCookieName)
		if err != nil {
			http.Redirect(w, r, "/terminal/login", http.StatusFound)
			return
		}

		restaurantID, err := strconv.Atoi(cookie.Value)
		if err != nil {
			clearSessionCookie(w)
			http.Redirect(w, r, "/terminal/login", http.StatusFound)
			return
		}

		ctx := context.WithValue(r.Context(), restaurantIDKey, int32(restaurantID))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func restaurantIDFromContext(ctx context.Context) int32 {
	id, _ := ctx.Value(restaurantIDKey).(int32)
	return id
}

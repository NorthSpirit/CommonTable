// internal/httpserver/terminal_routes.go
package httpserver

import (
	"CommonTable/internal/handlers/terminal"

	"github.com/go-chi/chi/v5"
)

type TerminalHandlers struct {
	Handler *terminal.Handler
}

func RegisterTerminalRoutes(r chi.Router, h TerminalHandlers) {
	r.Route("/terminal", func(r chi.Router) {
		r.Get("/login", h.Handler.LoginPage)
		r.Post("/login", h.Handler.Login)
		r.Get("/logout", h.Handler.Logout)

		r.Group(func(r chi.Router) {
			r.Use(terminal.RequireTerminalSession)
			r.Get("/", h.Handler.Board)
			r.Post("/orders/{orderID}/simulate", h.Handler.Simulate)
		})
	})
}

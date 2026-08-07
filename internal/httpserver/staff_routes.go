package httpserver

import (
	"CommonTable/internal/handlers/staff"
	"CommonTable/internal/services/authservices"

	"github.com/go-chi/chi/v5"
)

type StaffHandlers struct {
	Auth    *staff.AuthHandler
	Profile *staff.ProfileHandler
	Chat    *staff.ChatHandler
}

func RegisterStaffRoutes(r chi.Router, h StaffHandlers, authService *authservices.AuthService) {
	r.Route("/staff", func(r chi.Router) {
		r.Get("/login", h.Auth.RestaurantLoginPage)
		r.Post("/login", h.Auth.RestaurantLogin)
		r.Get("/login/employee", h.Auth.EmployeeLoginPage)
		r.Post("/login/employee", h.Auth.EmployeeLogin)
		r.Get("/logout", h.Auth.Logout)

		r.Group(func(r chi.Router) {
			r.Use(staff.RequireStaff(authService))
			r.Get("/profile", h.Profile.Edit)
			r.Post("/profile", h.Profile.Update)
			r.Get("/chat", h.Chat.Page)
			r.Post("/chat/send", h.Chat.Send)
			r.Get("/chat/messages", h.Chat.Poll)
		})
	})
}

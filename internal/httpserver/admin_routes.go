package httpserver

import (
	"CommonTable/internal/handlers/admin"
	"CommonTable/internal/services/adminservices"

	"github.com/go-chi/chi/v5"
)

type AdminHandlers struct {
	Auth      *admin.AuthHandler
	Dashboard *admin.DashboardHandler
	Homepage  *admin.HomepageHandler
}

func RegisterAdminRoutes(r chi.Router, h AdminHandlers, adminService *adminservices.AdminService) {
	r.Route("/admin", func(r chi.Router) {
		r.Get("/login", h.Auth.LoginPage)
		r.Post("/login", h.Auth.Login)
		r.Get("/logout", h.Auth.Logout)

		r.Group(func(r chi.Router) {
			r.Use(admin.RequireAdmin(adminService))
			r.Get("/", h.Dashboard.Dashboard)
			r.Get("/homepage", h.Homepage.Edit)
			r.Post("/homepage", h.Homepage.Update)
		})
	})
}

package container

import (
	"CommonTable/internal/handlers/admin"
	"CommonTable/internal/httpserver"
)

func NewAdminHandlers(s *Services) httpserver.AdminHandlers {
	return httpserver.AdminHandlers{
		Auth:      admin.NewAuthHandler(s.Admin, s.JWTTokenTTL),
		Dashboard: admin.NewDashboardHandler(),
		Homepage:  admin.NewHomepageHandler(s.Homepage, s.UploadsDir),
	}
}

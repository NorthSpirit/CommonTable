package container

import (
	"CommonTable/internal/handlers/staff"
	"CommonTable/internal/httpserver"
)

func NewStaffHandlers(s *Services) httpserver.StaffHandlers {
	return httpserver.StaffHandlers{
		Auth:    staff.NewAuthHandler(s.Auth, s.Restaurant, s.JWTTokenTTL),
		Profile: staff.NewProfileHandler(s.Employee),
		Chat:    staff.NewChatHandler(s.ChatMessage, s.Restaurant),
	}
}

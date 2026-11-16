package container

import "CommonTable/internal/handlers/admin"

func NewAdminApiKeyHandler(s *Services, ext *ExternalOrderServices) *admin.ApiKeyHandler {
	return admin.NewApiKeyHandler(ext.ApiKey, s.Restaurant)
}

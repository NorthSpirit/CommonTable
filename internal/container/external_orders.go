package container

import (
	"CommonTable/internal/handlers/admin"
	"CommonTable/internal/handlers/externalapi"
	"CommonTable/internal/handlers/manager"
	"CommonTable/internal/httpserver"
	"CommonTable/internal/repository"
	"CommonTable/internal/services/externalorderservices"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ExternalOrderServices struct {
	ApiKey *externalorderservices.ApiKeyService
	Order  *externalorderservices.ExternalOrderService
}

func NewExternalOrderServices(queries *repository.Queries, pool *pgxpool.Pool) *ExternalOrderServices {
	return &ExternalOrderServices{
		ApiKey: externalorderservices.NewApiKeyService(queries),
		Order:  externalorderservices.NewExternalOrderService(queries, pool),
	}
}

func NewExternalAPIHandlers(s *ExternalOrderServices) httpserver.ExternalAPIHandlers {
	return httpserver.ExternalAPIHandlers{
		Menu:  externalapi.NewMenuHandler(s.Order),
		Order: externalapi.NewOrderHandler(s.Order),
	}
}

func NewManagerOnlineOrderHandler(s *Services, ext *ExternalOrderServices) *manager.OnlineOrderHandler {
	return manager.NewOnlineOrderHandler(ext.Order, s.Restaurant)
}

func NewAdminOnlineOrderHandler(s *Services, ext *ExternalOrderServices) *admin.OnlineOrderHandler {
	return admin.NewOnlineOrderHandler(ext.Order, s.Employee, s.Restaurant)
}

func NewManagerApiKeyHandler(s *Services, ext *ExternalOrderServices) *manager.ApiKeyHandler {
	return manager.NewApiKeyHandler(ext.ApiKey, s.Restaurant)
}

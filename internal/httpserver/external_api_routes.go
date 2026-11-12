package httpserver

import (
	"CommonTable/internal/handlers/externalapi"
	"CommonTable/internal/services/externalorderservices"

	"github.com/go-chi/chi/v5"
)

type ExternalAPIHandlers struct {
	Menu  *externalapi.MenuHandler
	Order *externalapi.OrderHandler
}

func RegisterExternalAPIRoutes(r chi.Router, h ExternalAPIHandlers, keys *externalorderservices.ApiKeyService) {
	r.Route("/api/external", func(r chi.Router) {
		r.Use(externalapi.RequireApiKey(keys))
		r.Get("/items", h.Menu.ListItems)
		r.Post("/orders", h.Order.Create)
	})
}

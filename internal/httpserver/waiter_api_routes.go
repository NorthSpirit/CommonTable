package httpserver

import (
	"CommonTable/internal/handlers/waiterapi"
	"CommonTable/internal/services/authservices"

	"github.com/go-chi/chi/v5"
)

type WaiterAPIHandlers struct {
	Auth       *waiterapi.AuthHandler
	Me         *waiterapi.MeHandler
	Menu       *waiterapi.MenuHandler
	QuickItems *waiterapi.QuickItemsHandler
}

func RegisterWaiterAPIRoutes(r chi.Router, h WaiterAPIHandlers, authService *authservices.AuthService) {
	r.Route("/api/waiter", func(r chi.Router) {
		r.Post("/login", h.Auth.RestaurantLogin)
		r.Post("/login/employee", h.Auth.EmployeeLogin)

		r.Group(func(r chi.Router) {
			r.Use(waiterapi.RequireWaiterOrManager(authService))
			r.Get("/me", h.Me.Me)
			r.Get("/menus", h.Menu.ListMenus)
			r.Get("/menus/{menuID}/items", h.Menu.ListMenuItems)
			r.Get("/quick-items", h.QuickItems.List)
			r.Post("/quick-items", h.QuickItems.Add)
			r.Delete("/quick-items/{itemID}", h.QuickItems.Remove)
			r.Patch("/quick-items/{itemID}", h.QuickItems.Reorder)
		})
	})
}

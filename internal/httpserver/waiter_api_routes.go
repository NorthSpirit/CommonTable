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
	Chat       *waiterapi.ChatHandler
	Order      *waiterapi.OrderHandler
	Table      *waiterapi.TableHandler
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
			r.Get("/chat/messages", h.Chat.List)
			r.Get("/chat/messages/after/{id}", h.Chat.Poll)
			r.Post("/chat/send", h.Chat.Send)
			r.Get("/items/{itemID}", h.Menu.ItemDetail)
			r.Post("/orders", h.Order.Create)
			r.Get("/orders", h.Order.ListMyOrders)
			r.Delete("/orders/{orderID}/items/{itemOrderID}", h.Order.RemoveItem)
			r.Patch("/orders/{orderID}/items/{itemOrderID}", h.Order.UpdateItemQuantity)
			r.Post("/orders/{orderID}/cancel", h.Order.CancelOrder)
			r.Get("/tables", h.Table.List)
			r.Get("/search", h.Menu.Search)
		})
	})
}

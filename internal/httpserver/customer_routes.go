package httpserver

import (
	"CommonTable/internal/handlers/customer"

	"github.com/go-chi/chi/v5"
)

type CustomerHandlers struct {
	Home       *customer.HomeHandler
	Restaurant *customer.RestaurantHandler
	Menu       *customer.MenuHandler
	Item       *customer.ItemHandler
}

func RegisterCustomerRoutes(r chi.Router, h CustomerHandlers) {
	r.Get("/", h.Home.Home)
	r.Get("/restaurants", h.Restaurant.List)
	r.Get("/restaurants/{id}", h.Restaurant.Menus)
	r.Get("/menus/{id}", h.Menu.Items)
	r.Get("/items/{id}", h.Item.Detail)
}

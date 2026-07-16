package httpserver

import (
	"net/http"

	"CommonTable/internal/handlers/customer"
	"CommonTable/internal/templates/errorpages"

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

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		locale := customer.GetLocale(r)
		w.WriteHeader(http.StatusNotFound)
		_ = errorpages.NotFoundPage(locale).Render(r.Context(), w)
	})
}

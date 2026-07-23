package httpserver

import (
	"CommonTable/internal/handlers/chef"
	"CommonTable/internal/services/authservices"

	"github.com/go-chi/chi/v5"
)

type ChefHandlers struct {
	Auth *chef.AuthHandler
	Menu *chef.MenuHandler
}

func RegisterChefRoutes(r chi.Router, h ChefHandlers, authService *authservices.AuthService) {
	r.Route("/chef", func(r chi.Router) {
		r.Get("/login", h.Auth.RestaurantLoginPage)
		r.Post("/login", h.Auth.RestaurantLogin)
		r.Get("/login/employee", h.Auth.EmployeeLoginPage)
		r.Post("/login/employee", h.Auth.EmployeeLogin)
		r.Get("/logout", h.Auth.Logout)

		r.Group(func(r chi.Router) {
			r.Use(chef.RequireChef(authService))
			r.Get("/menus", h.Menu.List)
			r.Get("/menus/{menuID}", h.Menu.MenuDetail)
			r.Get("/items/{itemID}", h.Menu.ItemDetail)
		})
	})
}

package httpserver

import (
	"net/http"

	"CommonTable/internal/handlers/chef"
	"CommonTable/internal/services/authservices"

	"github.com/go-chi/chi/v5"
)

type ChefHandlers struct {
	Auth    *chef.AuthHandler
	Menu    *chef.MenuHandler
	Profile *chef.ProfileHandler
	Chat    *chef.ChatHandler
}

func RegisterChefRoutes(r chi.Router, h ChefHandlers, authService *authservices.AuthService) {
	r.Route("/chef", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/chef/login", http.StatusFound)
		})

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
			r.Post("/items/{itemID}/instructions", h.Menu.UpdateInstructions)
			r.Get("/profile", h.Profile.Edit)
			r.Post("/profile", h.Profile.Update)
			r.Get("/chat", h.Chat.Page)
			r.Post("/chat/send", h.Chat.Send)
			r.Get("/chat/messages", h.Chat.Poll)
		})
	})
}

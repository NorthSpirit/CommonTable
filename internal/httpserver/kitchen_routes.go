package httpserver

import (
	"net/http"

	"CommonTable/internal/handlers/kitchen"
	"CommonTable/internal/services/authservices"

	"github.com/go-chi/chi/v5"
)

type KitchenHandlers struct {
	Auth   *kitchen.AuthHandler
	Orders *kitchen.OrdersHandler
}

func RegisterKitchenRoutes(r chi.Router, h KitchenHandlers, authService *authservices.AuthService) {
	r.Route("/kitchen", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/kitchen/login", http.StatusFound)
		})

		r.Get("/login", h.Auth.RestaurantLoginPage)
		r.Post("/login", h.Auth.RestaurantLogin)
		r.Get("/login/employee", h.Auth.EmployeeLoginPage)
		r.Post("/login/employee", h.Auth.EmployeeLogin)
		r.Get("/logout", h.Auth.Logout)

		r.Group(func(r chi.Router) {
			r.Use(kitchen.RequireKitchenStaff(authService))
			r.Get("/orders", h.Orders.Board)
			r.Get("/orders/poll", h.Orders.Poll)
			r.Post("/orders/status", h.Orders.UpdateStatus)
		})
	})
}

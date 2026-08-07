package httpserver

import (
	"CommonTable/internal/handlers/manager"
	"CommonTable/internal/services/authservices"

	"github.com/go-chi/chi/v5"
)

type ManagerHandlers struct {
	Auth      *manager.AuthHandler
	Dashboard *manager.DashboardHandler
	Employee  *manager.EmployeeHandler
	Chat      *manager.ChatHandler
	Profile   *manager.ProfileHandler
	Order     *manager.OrderHandler
	Menu      *manager.MenuHandler
	FoodItem  *manager.FoodItemHandler
}

func RegisterManagerRoutes(r chi.Router, h ManagerHandlers, authService *authservices.AuthService) {
	r.Route("/manager", func(r chi.Router) {
		r.Get("/login", h.Auth.RestaurantLoginPage)
		r.Post("/login", h.Auth.RestaurantLogin)
		r.Get("/login/employee", h.Auth.EmployeeLoginPage)
		r.Post("/login/employee", h.Auth.EmployeeLogin)
		r.Get("/logout", h.Auth.Logout)

		r.Group(func(r chi.Router) {
			r.Use(manager.RequireManager(authService))
			r.Get("/", h.Dashboard.Dashboard)

			r.Route("/employees", func(r chi.Router) {
				r.Get("/", h.Employee.List)
				r.Get("/new", h.Employee.NewForm)
				r.Post("/new", h.Employee.Create)
				r.Get("/{employeeID}/edit", h.Employee.EditForm)
				r.Post("/{employeeID}/edit", h.Employee.Update)
				r.Post("/{employeeID}/toggle-active", h.Employee.ToggleActive)
			})

			r.Get("/chat", h.Chat.Page)
			r.Post("/chat/send", h.Chat.Send)
			r.Get("/chat/messages", h.Chat.Poll)

			r.Get("/profile", h.Profile.Edit)
			r.Post("/profile", h.Profile.Update)

			r.Route("/orders", func(r chi.Router) {
				r.Get("/", h.Order.List)
				r.Get("/new", h.Order.NewForm)
				r.Post("/new", h.Order.Create)
				r.Get("/{orderID}", h.Order.Detail)
				r.Post("/{orderID}", h.Order.Update)
				r.Post("/{orderID}/payments/{paymentID}", h.Order.UpdatePaymentStatus)
				r.Post("/{orderID}/items/{itemOrderID}/delete", h.Order.RemoveItem)
				r.Post("/{orderID}/items/{itemOrderID}/edit", h.Order.UpdateItemQuantity)
				r.Post("/{orderID}/payments/new", h.Order.CreatePayment)
				r.Post("/{orderID}/payments/quick-pay", h.Order.QuickPayCash)
			})

			r.Route("/menus", func(r chi.Router) {
				r.Get("/", h.Menu.List)
				r.Get("/new", h.Menu.NewForm)
				r.Post("/new", h.Menu.Create)
				r.Get("/{menuID}", h.Menu.Detail)
				r.Get("/{menuID}/edit", h.Menu.EditForm)
				r.Post("/{menuID}/edit", h.Menu.Update)
				r.Post("/{menuID}/delete", h.Menu.Delete)

				r.Route("/{menuID}/items", func(r chi.Router) {
					r.Get("/new", h.FoodItem.NewForm)
					r.Post("/new", h.FoodItem.Create)
					r.Get("/{itemID}/edit", h.FoodItem.EditForm)
					r.Post("/{itemID}/edit", h.FoodItem.Update)
					r.Post("/{itemID}/delete", h.FoodItem.Delete)
					r.Post("/{itemID}/restore", h.FoodItem.Restore)
					r.Post("/{itemID}/ingredients/add", h.FoodItem.AddIngredient)
					r.Post("/{itemID}/ingredients/create", h.FoodItem.CreateIngredient)
					r.Post("/{itemID}/ingredients/{ingredientID}/remove", h.FoodItem.RemoveIngredient)
				})
			})
		})
	})
}

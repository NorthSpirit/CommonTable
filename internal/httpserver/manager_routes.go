package httpserver

import (
	"CommonTable/internal/handlers/manager"
	"CommonTable/internal/services/authservices"

	"github.com/go-chi/chi/v5"
)

// Manager handlers is a toolbox for handlers the manager needs
// They are bundled together because they should be passed down all together and they are needed to be passed anyhow
type ManagerHandlers struct {
	Auth      *manager.AuthHandler
	Dashboard *manager.DashboardHandler
	Employee  *manager.EmployeeHandler
	Chat      *manager.ChatHandler
	Profile   *manager.ProfileHandler
	MyProfile *manager.MyProfileHandler
	Order     *manager.OrderHandler
	Menu      *manager.MenuHandler
	FoodItem  *manager.FoodItemHandler
}

// RegisterManagerRoutes attaches all /manager/... URLs to the main router `r`.
// Auth service is only used as a bouncer at the gates.
// Without chi all of this would be fairly more verbose, like full URL each time:
// POST /manager/orders/{orderID}/payments/{paymentID} for example
func RegisterManagerRoutes(r chi.Router, h ManagerHandlers, authService *authservices.AuthService) {
	// Obviously, everything below lives under /manager prefix, since it's a manager module
	r.Route("/manager", func(r chi.Router) {

		// Obviously login pages doesn't need auth yet
		// One is for the restaurant and the next one is for the person account
		r.Get("/login", h.Auth.RestaurantLoginPage)
		r.Post("/login", h.Auth.RestaurantLogin)
		r.Get("/login/employee", h.Auth.EmployeeLoginPage)
		r.Post("/login/employee", h.Auth.EmployeeLogin)
		r.Get("/logout", h.Auth.Logout)

		// Everything after login does.
		// So everything in this group does require auth.
		r.Group(func(r chi.Router) {
			// Runs before everything else in this group: checks the session cookie and that the role is "manager".
			// If it fails, the user is redirected to /manager/login (based on manager/middleware.go)
			r.Use(manager.RequireManager(authService))
			r.Get("/", h.Dashboard.Dashboard)

			// /manager/employees/...
			r.Route("/employees", func(r chi.Router) {
				r.Get("/", h.Employee.List)
				r.Get("/new", h.Employee.NewForm)
				r.Post("/new", h.Employee.Create)
				r.Get("/{employeeID}/edit", h.Employee.EditForm)
				r.Post("/{employeeID}/edit", h.Employee.Update)
				r.Post("/{employeeID}/toggle-active", h.Employee.ToggleActive)
			})

			// Chat: page itself, sending a message, and "any new messages?"
			r.Get("/chat", h.Chat.Page)
			r.Post("/chat/send", h.Chat.Send)
			r.Get("/chat/messages", h.Chat.Poll)

			// Restaurant profile: view/edit form, then save
			r.Get("/profile", h.Profile.Edit)
			r.Post("/profile", h.Profile.Update)

			// The logged-in user's personal profile: view/edit form, then save
			r.Get("/my-profile", h.MyProfile.Edit)
			r.Post("/my-profile", h.MyProfile.Update)

			// /manager/orders/...
			r.Route("/orders", func(r chi.Router) {
				r.Get("/", h.Order.List)             // all orders
				r.Get("/new", h.Order.NewForm)       // "create order" form
				r.Post("/new", h.Order.Create)       // save new order
				r.Get("/{orderID}", h.Order.Detail)  // view one order
				r.Post("/{orderID}", h.Order.Update) // save changes to that order

				// Payments on an order
				r.Post("/{orderID}/payments/{paymentID}", h.Order.UpdatePaymentStatus) // e.g. mark paid
				r.Post("/{orderID}/payments/new", h.Order.CreatePayment)               // add a payment
				r.Post("/{orderID}/payments/quick-pay", h.Order.QuickPayCash)          // one-click cash payment

				// Items (dishes) inside an order
				r.Post("/{orderID}/items/{itemOrderID}/delete", h.Order.RemoveItem)       // remove a dish
				r.Post("/{orderID}/items/{itemOrderID}/edit", h.Order.UpdateItemQuantity) // change quantity
				r.Post("/{orderID}/items/add", h.Order.AddItem)                           // add a dish
			})

			// /manager/menus/...
			r.Route("/menus", func(r chi.Router) {
				r.Get("/", h.Menu.List)
				r.Get("/new", h.Menu.NewForm)
				r.Post("/new", h.Menu.Create)
				r.Get("/{menuID}", h.Menu.Detail)
				r.Get("/{menuID}/edit", h.Menu.EditForm)
				r.Post("/{menuID}/edit", h.Menu.Update)
				r.Post("/{menuID}/delete", h.Menu.Delete)

				// Dishes inside a specific menu: /manager/menus/{menuID}/items/...
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

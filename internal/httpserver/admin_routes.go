package httpserver

import (
	"CommonTable/internal/handlers/admin"
	"CommonTable/internal/services/adminservices"

	"github.com/go-chi/chi/v5"
)

type AdminHandlers struct {
	Auth       *admin.AuthHandler
	Dashboard  *admin.DashboardHandler
	Homepage   *admin.HomepageHandler
	Restaurant *admin.RestaurantHandler
	Menu       *admin.MenuHandler
	FoodItem   *admin.FoodItemHandler
	Ingredient *admin.IngredientHandler
	Allergen   *admin.AllergenHandler
	Employee   *admin.EmployeeHandler
	Table      *admin.TableHandler
}

func RegisterAdminRoutes(r chi.Router, h AdminHandlers, adminService *adminservices.AdminService) {
	r.Route("/admin", func(r chi.Router) {
		r.Get("/login", h.Auth.LoginPage)
		r.Post("/login", h.Auth.Login)
		r.Get("/logout", h.Auth.Logout)

		r.Group(func(r chi.Router) {
			r.Use(admin.RequireAdmin(adminService))
			r.Get("/", h.Dashboard.Dashboard)
			r.Get("/homepage", h.Homepage.Edit)
			r.Post("/homepage", h.Homepage.Update)

			r.Route("/restaurants", func(r chi.Router) {
				r.Get("/", h.Restaurant.List)
				r.Get("/new", h.Restaurant.NewForm)
				r.Post("/new", h.Restaurant.Create)
				r.Get("/{id}", h.Restaurant.Detail)
				r.Get("/{id}/edit", h.Restaurant.EditForm)
				r.Post("/{id}/edit", h.Restaurant.Update)
				r.Get("/{id}/delete", h.Restaurant.DeleteConfirm)
				r.Post("/{id}/delete", h.Restaurant.Delete)
				r.Post("/{id}/restore", h.Restaurant.Restore)

				r.Route("/{id}/menus", func(r chi.Router) {
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
					})
				})

				r.Route("/{id}/employees", func(r chi.Router) {
					r.Get("/", h.Employee.List)
					r.Get("/new", h.Employee.NewForm)
					r.Post("/new", h.Employee.Create)
					r.Get("/{employeeID}/edit", h.Employee.EditForm)
					r.Post("/{employeeID}/edit", h.Employee.Update)
					r.Post("/{employeeID}/toggle-active", h.Employee.ToggleActive)
				})
			})

			r.Route("/ingredients", func(r chi.Router) {
				r.Get("/", h.Ingredient.List)
				r.Get("/new", h.Ingredient.NewForm)
				r.Post("/new", h.Ingredient.Create)
				r.Get("/{id}/edit", h.Ingredient.EditForm)
				r.Post("/{id}/edit", h.Ingredient.Update)
				r.Post("/{id}/delete", h.Ingredient.Delete)
			})

			r.Route("/allergens", func(r chi.Router) {
				r.Get("/", h.Allergen.List)
				r.Get("/new", h.Allergen.NewForm)
				r.Post("/new", h.Allergen.Create)
				r.Get("/{id}/edit", h.Allergen.EditForm)
				r.Post("/{id}/edit", h.Allergen.Update)
				r.Post("/{id}/delete", h.Allergen.Delete)
			})

			r.Route("/tables", func(r chi.Router) {
				r.Get("/", h.Table.List)
				r.Get("/new", h.Table.NewForm)
				r.Post("/new", h.Table.Create)
				r.Get("/{id}/edit", h.Table.EditForm)
				r.Post("/{id}/edit", h.Table.Update)
			})
		})
	})
}

package container

import (
	"CommonTable/internal/handlers/admin"
	"CommonTable/internal/httpserver"
)

func NewAdminHandlers(s *Services) httpserver.AdminHandlers {
	return httpserver.AdminHandlers{
		Auth:       admin.NewAuthHandler(s.Admin, s.JWTTokenTTL),
		Dashboard:  admin.NewDashboardHandler(),
		Homepage:   admin.NewHomepageHandler(s.Homepage, s.UploadsDir),
		Restaurant: admin.NewRestaurantHandler(s.Restaurant, s.Admin, s.UploadsDir),
		Menu:       admin.NewMenuHandler(s.FoodMenu, s.Restaurant, s.ItemMenu, s.UploadsDir),
		FoodItem:   admin.NewFoodItemHandler(s.FoodItem, s.FoodMenu, s.FoodAllergen, s.ItemAllergen, s.ItemMenu, s.UploadsDir),
		Ingredient: admin.NewIngredientHandler(s.FoodIngredient),
		Allergen:   admin.NewAllergenHandler(s.FoodAllergen),
		Employee:   admin.NewEmployeeHandler(s.Employee, s.Restaurant),
		Table:      admin.NewTableHandler(s.Table),
	}
}

package container

import (
	"CommonTable/internal/handlers/admin"
	"CommonTable/internal/httpserver"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewAdminHandlers(s *Services, pool *pgxpool.Pool) httpserver.AdminHandlers {
	return httpserver.AdminHandlers{
		Auth:           admin.NewAuthHandler(s.Admin, s.JWTTokenTTL),
		Dashboard:      admin.NewDashboardHandler(s.Restaurant, s.FoodMenu, s.FoodItem, s.Employee, s.FoodIngredient, s.FoodAllergen, s.Table),
		Homepage:       admin.NewHomepageHandler(s.Homepage, s.UploadsDir),
		Restaurant:     admin.NewRestaurantHandler(s.Restaurant, s.Admin, s.UploadsDir),
		Menu:           admin.NewMenuHandler(s.FoodMenu, s.Restaurant, s.ItemMenu, s.UploadsDir),
		FoodItem:       admin.NewFoodItemHandler(s.FoodItem, s.FoodMenu, s.FoodAllergen, s.FoodIngredient, s.ItemAllergen, s.ItemMenu, s.IngredientItem, s.UploadsDir),
		Allergen:       admin.NewAllergenHandler(s.FoodAllergen, s.ItemAllergen),
		Ingredient:     admin.NewIngredientHandler(s.FoodIngredient, s.IngredientItem),
		Employee:       admin.NewEmployeeHandler(s.Employee, s.Restaurant),
		Table:          admin.NewTableHandler(s.Table, s.Restaurant),
		Order:          admin.NewOrderHandler(pool, s.FoodOrder, s.ItemOrder, s.Payment, s.Table, s.Employee, s.Restaurant, s.FoodMenu, s.ItemMenu, s.FoodItem),
		Chat:           admin.NewChatHandler(s.ChatMessage, s.Restaurant),
		ChangePassword: admin.NewChangePasswordHandler(s.Admin),
		TableStatus:    admin.NewTableStatusHandler(s.Table, s.Restaurant),
		Export:         admin.NewExportHandler(s.Restaurant, s.FoodMenu, s.ItemMenu, s.ItemAllergen, s.IngredientItem),
	}
}

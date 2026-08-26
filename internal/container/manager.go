package container

import (
	"CommonTable/internal/handlers/manager"
	"CommonTable/internal/httpserver"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewManagerHandlers(s *Services, pool *pgxpool.Pool) httpserver.ManagerHandlers {
	return httpserver.ManagerHandlers{
		Auth:      manager.NewAuthHandler(s.Auth, s.Restaurant, s.JWTTokenTTL),
		Dashboard: manager.NewDashboardHandler(s.Restaurant, s.FoodMenu, s.FoodItem, s.Employee, s.FoodOrder),
		Employee:  manager.NewEmployeeHandler(s.Employee, s.Restaurant),
		Chat:      manager.NewChatHandler(s.ChatMessage, s.Restaurant),
		Profile:   manager.NewProfileHandler(s.Restaurant, s.UploadsDir),
		Order:     manager.NewOrderHandler(pool, s.FoodOrder, s.ItemOrder, s.Payment, s.Table, s.Employee, s.Restaurant, s.FoodMenu, s.ItemMenu, s.FoodItem),
		Menu:      manager.NewMenuHandler(s.FoodMenu, s.ItemMenu, s.Restaurant, s.UploadsDir),
		FoodItem:  manager.NewFoodItemHandler(s.FoodItem, s.FoodMenu, s.FoodAllergen, s.FoodIngredient, s.ItemAllergen, s.ItemMenu, s.IngredientItem, s.UploadsDir),
		MyProfile: manager.NewMyProfileHandler(s.Employee),
	}
}

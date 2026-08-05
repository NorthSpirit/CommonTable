package container

import (
	"CommonTable/internal/handlers/waiterapi"
	"CommonTable/internal/httpserver"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewWaiterAPIHandlers(s *Services, pool *pgxpool.Pool) httpserver.WaiterAPIHandlers {
	return httpserver.WaiterAPIHandlers{
		Auth:       waiterapi.NewAuthHandler(s.Auth),
		Me:         waiterapi.NewMeHandler(),
		Menu:       waiterapi.NewMenuHandler(s.FoodMenu, s.FoodItem, s.ItemAllergen, s.IngredientItem),
		QuickItems: waiterapi.NewQuickItemsHandler(s.WaiterQuickItem),
		Chat:       waiterapi.NewChatHandler(s.ChatMessage),
		Order:      waiterapi.NewOrderHandler(pool, s.FoodOrder, s.ItemOrder, s.FoodItem),
		Table:      waiterapi.NewTableHandler(s.Table),
	}
}

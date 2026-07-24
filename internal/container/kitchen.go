package container

import (
	"CommonTable/internal/handlers/kitchen"
	"CommonTable/internal/httpserver"
)

func NewKitchenHandlers(s *Services) httpserver.KitchenHandlers {
	return httpserver.KitchenHandlers{
		Auth:   kitchen.NewAuthHandler(s.Auth, s.Restaurant, s.JWTTokenTTL),
		Orders: kitchen.NewOrdersHandler(s.FoodOrder, s.ItemOrder),
	}
}

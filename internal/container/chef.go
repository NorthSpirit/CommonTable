package container

import (
	"CommonTable/internal/handlers/chef"
	"CommonTable/internal/httpserver"
)

func NewChefHandlers(s *Services) httpserver.ChefHandlers {
	return httpserver.ChefHandlers{
		Auth:    chef.NewAuthHandler(s.Auth, s.Restaurant, s.JWTTokenTTL),
		Menu:    chef.NewMenuHandler(s.FoodMenu, s.FoodItem, s.ItemMenu, s.ItemAllergen, s.IngredientItem, s.Restaurant),
		Profile: chef.NewProfileHandler(s.Employee),
		Chat:    chef.NewChatHandler(s.ChatMessage, s.Restaurant),
	}
}

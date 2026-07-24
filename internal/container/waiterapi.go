package container

import (
	"CommonTable/internal/handlers/waiterapi"
	"CommonTable/internal/httpserver"
)

func NewWaiterAPIHandlers(s *Services) httpserver.WaiterAPIHandlers {
	return httpserver.WaiterAPIHandlers{
		Auth:       waiterapi.NewAuthHandler(s.Auth),
		Me:         waiterapi.NewMeHandler(),
		Menu:       waiterapi.NewMenuHandler(s.FoodMenu, s.FoodItem),
		QuickItems: waiterapi.NewQuickItemsHandler(s.WaiterQuickItem),
	}
}

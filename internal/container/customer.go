package container

import (
	"CommonTable/internal/handlers/customer"
	"CommonTable/internal/httpserver"
)

func NewCustomerHandlers(s *Services) httpserver.CustomerHandlers {
	return httpserver.CustomerHandlers{
		Home:       customer.NewHomeHandler(),
		Restaurant: customer.NewRestaurantHandler(s.Restaurant, s.FoodMenu),
		Menu:       customer.NewMenuHandler(s.FoodMenu, s.FoodItem),
		Item:       customer.NewItemHandler(s.FoodItem, s.ItemAllergen),
	}
}

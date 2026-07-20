package container

import (
	"CommonTable/internal/handlers/customer"
	"CommonTable/internal/httpserver"
)

func NewCustomerHandlers(s *Services) httpserver.CustomerHandlers {
	return httpserver.CustomerHandlers{
		Home:       customer.NewHomeHandler(s.Homepage, s.Restaurant),
		Restaurant: customer.NewRestaurantHandler(s.Restaurant, s.FoodMenu, s.FoodItem),
		Menu:       customer.NewMenuHandler(s.FoodMenu, s.FoodItem, s.Restaurant),
		Item:       customer.NewItemHandler(s.FoodItem, s.ItemAllergen, s.Restaurant),
	}
}

package container

import (
	"CommonTable/internal/repository"
	"CommonTable/internal/services/adminservices"
	"CommonTable/internal/services/authservices"
	"CommonTable/internal/services/foodrelationshipservices"
	"CommonTable/internal/services/foodservices"
	"CommonTable/internal/services/restaurantservices"

	"time"
)

// Services holds every service the application needs, constructed once
// at startup and passed down to whichever handler builders need them.
type Services struct {
	Auth  *authservices.AuthService
	Admin *adminservices.AdminService

	Restaurant   *restaurantservices.RestaurantService
	FoodMenu     *foodservices.FoodMenuService
	FoodItem     *foodservices.FoodItemService
	ItemAllergen *foodrelationshipservices.ItemAllergenService

	// add more here as new domains get built —
	// Employee, Table, FoodOrder, Payment, ChatMessage, etc.
}

type ServicesConfig struct {
	AdminUsername     string
	AdminPasswordHash string
	JWTSecret         []byte
	JWTTokenTTL       time.Duration
}

func NewServices(queries *repository.Queries, cfg ServicesConfig) *Services {
	return &Services{
		Auth: authservices.NewAuthService(queries, cfg.JWTSecret, cfg.JWTTokenTTL),
		Admin: adminservices.NewAdminService(
			adminservices.AdminConfig{
				Username:     cfg.AdminUsername,
				PasswordHash: cfg.AdminPasswordHash,
			},
			cfg.JWTSecret,
			cfg.JWTTokenTTL,
		),

		Restaurant:   restaurantservices.NewRestaurantService(queries),
		FoodMenu:     foodservices.NewFoodMenuService(queries),
		FoodItem:     foodservices.NewFoodItemService(queries),
		ItemAllergen: foodrelationshipservices.NewItemAllergenService(queries),
	}
}

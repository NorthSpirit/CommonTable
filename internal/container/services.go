package container

import (
	"time"

	"CommonTable/internal/repository"
	"CommonTable/internal/services/adminservices"
	"CommonTable/internal/services/authservices"
	"CommonTable/internal/services/foodrelationshipservices"
	"CommonTable/internal/services/foodservices"
	"CommonTable/internal/services/homepageservices"
	"CommonTable/internal/services/restaurantservices"
)

// Services holds every service the application needs, constructed once
// at startup and passed down to whichever handler builders need them.
type Services struct {
	Auth  *authservices.AuthService
	Admin *adminservices.AdminService

	Restaurant   *restaurantservices.RestaurantService
	FoodMenu     *foodservices.FoodMenuService
	FoodItem     *foodservices.FoodItemService
	FoodAllergen *foodservices.FoodAllergenService
	ItemAllergen *foodrelationshipservices.ItemAllergenService
	Homepage     *homepageservices.HomepageService

	JWTTokenTTL time.Duration
	UploadsDir  string
}

type ServicesConfig struct {
	AdminUsername     string
	AdminPasswordHash string
	JWTSecret         []byte
	JWTTokenTTL       time.Duration
	UploadsDir        string
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
		FoodAllergen: foodservices.NewFoodAllergenService(queries),
		ItemAllergen: foodrelationshipservices.NewItemAllergenService(queries),
		Homepage:     homepageservices.NewHomepageService(queries),
		JWTTokenTTL:  cfg.JWTTokenTTL,
		UploadsDir:   cfg.UploadsDir,
	}
}

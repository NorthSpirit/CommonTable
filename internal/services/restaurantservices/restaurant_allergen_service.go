package restaurantservices

import (
	"context"
	"fmt"

	"CommonTable/internal/repository"
)

type RestaurantAllergenService struct {
	queries *repository.Queries
}

func NewRestaurantAllergenService(queries *repository.Queries) *RestaurantAllergenService {
	return &RestaurantAllergenService{queries: queries}
}

func (f *RestaurantAllergenService) ListAllergensForRestaurant(ctx context.Context, restaurantID int32) ([]repository.Allergen, error) {
	allergens, err := f.queries.ListAllergensForRestaurant(ctx, restaurantID)
	if err != nil {
		return nil, fmt.Errorf("list allergens for restaurant: %w", err)
	}

	return allergens, nil
}

func (f *RestaurantAllergenService) LinkRestaurantAllergen(ctx context.Context, restaurantID int32, allergenID int32) error {
	err := f.queries.LinkRestaurantAllergen(ctx, repository.LinkRestaurantAllergenParams{
		RestaurantID: restaurantID,
		AllergenID:   allergenID,
	})
	if err != nil {
		return fmt.Errorf("link restaurant allergen: %w", err)
	}

	return nil
}

func (f *RestaurantAllergenService) UnlinkRestaurantAllergen(ctx context.Context, restaurantID int32, allergenID int32) error {
	err := f.queries.UnlinkRestaurantAllergen(ctx, repository.UnlinkRestaurantAllergenParams{
		RestaurantID: restaurantID,
		AllergenID:   allergenID,
	})
	if err != nil {
		return fmt.Errorf("unlink restaurant allergen: %w", err)
	}

	return nil
}

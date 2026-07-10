package restaurantservices

import (
	"context"
	"fmt"

	"CommonTable/internal/repository"
)

type RestaurantIngredientService struct {
	queries *repository.Queries
}

func NewRestaurantIngredientService(queries *repository.Queries) *RestaurantIngredientService {
	return &RestaurantIngredientService{queries: queries}
}

func (f *RestaurantIngredientService) ListIngredientsForRestaurant(ctx context.Context, restaurantID int32) ([]repository.ListIngredientsForRestaurantRow, error) {
	ingredients, err := f.queries.ListIngredientsForRestaurant(ctx, restaurantID)
	if err != nil {
		return nil, fmt.Errorf("list ingredients for restaurant: %w", err)
	}

	return ingredients, nil
}

func (f *RestaurantIngredientService) LinkRestaurantIngredient(
	ctx context.Context,
	restaurantID int32,
	ingredientID int32,
	isShared bool,
) error {

	err := f.queries.LinkRestaurantIngredient(ctx, repository.LinkRestaurantIngredientParams{
		RestaurantID: restaurantID,
		IngredientID: ingredientID,
		IsShared:     isShared,
	})
	if err != nil {
		return fmt.Errorf("link restaurant ingredient: %w", err)
	}

	return nil
}

func (f *RestaurantIngredientService) UpdateRestaurantIngredientShared(
	ctx context.Context,
	restaurantID int32,
	ingredientID int32,
	isShared bool,
) error {

	err := f.queries.UpdateRestaurantIngredientShared(ctx, repository.UpdateRestaurantIngredientSharedParams{
		RestaurantID: restaurantID,
		IngredientID: ingredientID,
		IsShared:     isShared,
	})
	if err != nil {
		return fmt.Errorf("update restaurant ingredient shared: %w", err)
	}

	return nil
}

func (f *RestaurantIngredientService) UnlinkRestaurantIngredient(ctx context.Context, restaurantID int32, ingredientID int32) error {
	err := f.queries.UnlinkRestaurantIngredient(ctx, repository.UnlinkRestaurantIngredientParams{
		RestaurantID: restaurantID,
		IngredientID: ingredientID,
	})
	if err != nil {
		return fmt.Errorf("unlink restaurant ingredient: %w", err)
	}

	return nil
}

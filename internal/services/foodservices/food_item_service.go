package foodservices

import (
	"context"
	"fmt"

	"CommonTable/internal/helpers"
	"CommonTable/internal/repository"
)

type FoodItemService struct {
	queries *repository.Queries
}

func NewFoodItemService(queries *repository.Queries) *FoodItemService {
	return &FoodItemService{queries: queries}
}

func (f *FoodItemService) GetFoodItemByID(ctx context.Context, id int32) (repository.GetFoodItemRow, error) {
	item, err := f.queries.GetFoodItem(ctx, id)
	if err != nil {
		return repository.GetFoodItemRow{}, fmt.Errorf("get food item: %w", err)
	}

	return item, nil
}

func (f *FoodItemService) ListFoodItemsByRestaurant(ctx context.Context, restaurantID int32) ([]repository.ListFoodItemsByRestaurantRow, error) {
	items, err := f.queries.ListFoodItemsByRestaurant(ctx, restaurantID)
	if err != nil {
		return nil, fmt.Errorf("list food items by restaurant: %w", err)
	}

	return items, nil
}

func (f *FoodItemService) ListAvailableFoodItemsByRestaurant(ctx context.Context, restaurantID int32) ([]repository.ListAvailableFoodItemsByRestaurantRow, error) {
	items, err := f.queries.ListAvailableFoodItemsByRestaurant(ctx, restaurantID)
	if err != nil {
		return nil, fmt.Errorf("list available food items by restaurant: %w", err)
	}

	return items, nil
}

// ListMenuItemsForDisplay returns available items on a given menu, including
// price, for public-facing menu display (no authentication required).
func (f *FoodItemService) ListMenuItemsForDisplay(ctx context.Context, menuID int32) ([]repository.ListMenuItemsForDisplayRow, error) {
	items, err := f.queries.ListMenuItemsForDisplay(ctx, menuID)
	if err != nil {
		return nil, fmt.Errorf("list menu items for display: %w", err)
	}

	return items, nil
}

func (f *FoodItemService) InsertFoodItem(
	ctx context.Context,
	restaurantID int32,
	name string,
	description *string,
	isAvailable bool,
	imageURL *string,
	prepTimeMins *int32,
	instructions *string,
) (repository.CreateFoodItemRow, error) {

	item, err := f.queries.CreateFoodItem(ctx, repository.CreateFoodItemParams{
		RestaurantID: restaurantID,
		Name:         name,
		Description:  helpers.TextFromString(description),
		IsAvailable:  isAvailable,
		ImageUrl:     helpers.TextFromString(imageURL),
		PrepTimeMins: helpers.Int4FromInt32(prepTimeMins),
		Instructions: helpers.TextFromString(instructions),
	})
	if err != nil {
		return repository.CreateFoodItemRow{}, fmt.Errorf("create food item: %w", err)
	}

	return item, nil
}

func (f *FoodItemService) UpdateFoodItem(
	ctx context.Context,
	id int32,
	name string,
	description *string,
	isAvailable bool,
	imageURL *string,
	prepTimeMins *int32,
	instructions *string,
) error {

	err := f.queries.UpdateFoodItem(ctx, repository.UpdateFoodItemParams{
		ID:           id,
		Name:         name,
		Description:  helpers.TextFromString(description),
		IsAvailable:  isAvailable,
		ImageUrl:     helpers.TextFromString(imageURL),
		PrepTimeMins: helpers.Int4FromInt32(prepTimeMins),
		Instructions: helpers.TextFromString(instructions),
	})
	if err != nil {
		return fmt.Errorf("update food item: %w", err)
	}

	return nil
}

func (f *FoodItemService) SetFoodItemAvailability(ctx context.Context, id int32, isAvailable bool) error {
	err := f.queries.SetFoodItemAvailability(ctx, repository.SetFoodItemAvailabilityParams{
		ID:          id,
		IsAvailable: isAvailable,
	})
	if err != nil {
		return fmt.Errorf("set food item availability: %w", err)
	}

	return nil
}

func (f *FoodItemService) DeleteFoodItem(ctx context.Context, id int32) error {
	err := f.queries.SoftDeleteFoodItem(ctx, id)
	if err != nil {
		return fmt.Errorf("delete food item: %w", err)
	}

	return nil
}

func (f *FoodItemService) RestoreFoodItem(ctx context.Context, id int32) error {
	err := f.queries.RestoreFoodItem(ctx, id)
	if err != nil {
		return fmt.Errorf("restore food item: %w", err)
	}

	return nil
}

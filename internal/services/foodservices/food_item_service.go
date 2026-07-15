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

// GetFoodItemByID returns raw item data in both languages — for admin editing.
func (f *FoodItemService) GetFoodItemByID(ctx context.Context, id int32) (repository.GetFoodItemRow, error) {
	item, err := f.queries.GetFoodItem(ctx, id)
	if err != nil {
		return repository.GetFoodItemRow{}, fmt.Errorf("get food item: %w", err)
	}

	return item, nil
}

// GetFoodItemForDisplay returns locale-resolved item data — for customer-facing pages.
func (f *FoodItemService) GetFoodItemForDisplay(ctx context.Context, id int32, locale string) (repository.GetFoodItemForDisplayRow, error) {
	item, err := f.queries.GetFoodItemForDisplay(ctx, repository.GetFoodItemForDisplayParams{
		Locale: locale,
		ID:     id,
	})
	if err != nil {
		return repository.GetFoodItemForDisplayRow{}, fmt.Errorf("get food item for display: %w", err)
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

// ListMenuItemsForDisplay returns available items on a menu, locale-resolved,
// for public-facing menu display.
func (f *FoodItemService) ListMenuItemsForDisplay(ctx context.Context, menuID int32, locale string) ([]repository.ListMenuItemsForDisplayRow, error) {
	items, err := f.queries.ListMenuItemsForDisplay(ctx, repository.ListMenuItemsForDisplayParams{
		Locale: locale,
		MenuID: menuID,
	})
	if err != nil {
		return nil, fmt.Errorf("list menu items for display: %w", err)
	}

	return items, nil
}

func (f *FoodItemService) InsertFoodItem(
	ctx context.Context,
	restaurantID int32,
	name string,
	nameAlt *string,
	description *string,
	descriptionAlt *string,
	isAvailable bool,
	imageURL *string,
	prepTimeMins *int32,
	instructions *string,
	instructionsAlt *string,
) (repository.CreateFoodItemRow, error) {

	item, err := f.queries.CreateFoodItem(ctx, repository.CreateFoodItemParams{
		RestaurantID:    restaurantID,
		Name:            name,
		NameAlt:         helpers.TextFromString(nameAlt),
		Description:     helpers.TextFromString(description),
		DescriptionAlt:  helpers.TextFromString(descriptionAlt),
		IsAvailable:     isAvailable,
		ImageUrl:        helpers.TextFromString(imageURL),
		PrepTimeMins:    helpers.Int4FromInt32(prepTimeMins),
		Instructions:    helpers.TextFromString(instructions),
		InstructionsAlt: helpers.TextFromString(instructionsAlt),
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
	nameAlt *string,
	description *string,
	descriptionAlt *string,
	isAvailable bool,
	imageURL *string,
	prepTimeMins *int32,
	instructions *string,
	instructionsAlt *string,
) error {

	err := f.queries.UpdateFoodItem(ctx, repository.UpdateFoodItemParams{
		ID:              id,
		Name:            name,
		NameAlt:         helpers.TextFromString(nameAlt),
		Description:     helpers.TextFromString(description),
		DescriptionAlt:  helpers.TextFromString(descriptionAlt),
		IsAvailable:     isAvailable,
		ImageUrl:        helpers.TextFromString(imageURL),
		PrepTimeMins:    helpers.Int4FromInt32(prepTimeMins),
		Instructions:    helpers.TextFromString(instructions),
		InstructionsAlt: helpers.TextFromString(instructionsAlt),
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

// ListMostOrderedItemsByRestaurant returns the top N most-ordered available
// items for a restaurant, locale-resolved — used for the quick-nav.
func (f *FoodItemService) ListMostOrderedItemsByRestaurant(ctx context.Context, restaurantID int32, locale string, limit int32) ([]repository.ListMostOrderedItemsByRestaurantRow, error) {
	items, err := f.queries.ListMostOrderedItemsByRestaurant(ctx, repository.ListMostOrderedItemsByRestaurantParams{
		Locale:       locale,
		RestaurantID: restaurantID,
		ItemLimit:    limit,
	})
	if err != nil {
		return nil, fmt.Errorf("list most ordered items by restaurant: %w", err)
	}

	return items, nil
}

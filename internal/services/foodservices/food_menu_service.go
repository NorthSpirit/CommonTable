package foodservices

import (
	"context"
	"fmt"

	"CommonTable/internal/helpers"
	"CommonTable/internal/repository"
)

type FoodMenuService struct {
	queries *repository.Queries
}

func NewFoodMenuService(queries *repository.Queries) *FoodMenuService {
	return &FoodMenuService{queries: queries}
}

func (f *FoodMenuService) GetFoodMenuByID(ctx context.Context, id int32) (repository.FoodMenu, error) {
	menu, err := f.queries.GetFoodMenu(ctx, id)
	if err != nil {
		return repository.FoodMenu{}, fmt.Errorf("get food menu: %w", err)
	}

	return menu, nil
}

func (f *FoodMenuService) GetFoodMenuForDisplay(ctx context.Context, id int32, locale string) (repository.GetFoodMenuForDisplayRow, error) {
	menu, err := f.queries.GetFoodMenuForDisplay(ctx, repository.GetFoodMenuForDisplayParams{
		Locale: locale,
		ID:     id,
	})
	if err != nil {
		return repository.GetFoodMenuForDisplayRow{}, fmt.Errorf("get food menu for display: %w", err)
	}

	return menu, nil
}

func (f *FoodMenuService) ListFoodMenusByRestaurant(ctx context.Context, restaurantID int32) ([]repository.FoodMenu, error) {
	menus, err := f.queries.ListFoodMenusByRestaurant(ctx, restaurantID)
	if err != nil {
		return nil, fmt.Errorf("list food menus by restaurant: %w", err)
	}

	return menus, nil
}

func (f *FoodMenuService) ListFoodMenusByRestaurantForDisplay(ctx context.Context, restaurantID int32, locale string) ([]repository.ListFoodMenusByRestaurantForDisplayRow, error) {
	menus, err := f.queries.ListFoodMenusByRestaurantForDisplay(ctx, repository.ListFoodMenusByRestaurantForDisplayParams{
		Locale:       locale,
		RestaurantID: restaurantID,
	})
	if err != nil {
		return nil, fmt.Errorf("list food menus by restaurant for display: %w", err)
	}

	return menus, nil
}

func (f *FoodMenuService) InsertFoodMenu(
	ctx context.Context,
	restaurantID int32,
	name string,
	nameAlt *string,
	description *string,
	descriptionAlt *string,
	imageURL *string,
) (repository.FoodMenu, error) {

	menu, err := f.queries.CreateFoodMenu(ctx, repository.CreateFoodMenuParams{
		RestaurantID:   restaurantID,
		Name:           name,
		NameAlt:        helpers.TextFromString(nameAlt),
		Description:    helpers.TextFromString(description),
		DescriptionAlt: helpers.TextFromString(descriptionAlt),
		ImageUrl:       helpers.TextFromString(imageURL),
	})
	if err != nil {
		return repository.FoodMenu{}, fmt.Errorf("create food menu: %w", err)
	}

	return menu, nil
}

func (f *FoodMenuService) UpdateFoodMenu(
	ctx context.Context,
	id int32,
	name string,
	nameAlt *string,
	description *string,
	descriptionAlt *string,
	imageURL *string,
) error {

	err := f.queries.UpdateFoodMenu(ctx, repository.UpdateFoodMenuParams{
		ID:             id,
		Name:           name,
		NameAlt:        helpers.TextFromString(nameAlt),
		Description:    helpers.TextFromString(description),
		DescriptionAlt: helpers.TextFromString(descriptionAlt),
		ImageUrl:       helpers.TextFromString(imageURL),
	})
	if err != nil {
		return fmt.Errorf("update food menu: %w", err)
	}

	return nil
}

func (f *FoodMenuService) DeleteFoodMenu(ctx context.Context, id int32) error {
	err := f.queries.DeleteFoodMenu(ctx, id)
	if err != nil {
		return fmt.Errorf("delete food menu: %w", err)
	}

	return nil
}

package foodrelationshipservices

import (
	"context"
	"fmt"

	"CommonTable/internal/repository"
)

type ItemMenuService struct {
	queries *repository.Queries
}

func NewItemMenuService(queries *repository.Queries) *ItemMenuService {
	return &ItemMenuService{queries: queries}
}

func (f *ItemMenuService) ListMenusForItem(ctx context.Context, itemID int32) ([]repository.FoodMenu, error) {
	menus, err := f.queries.ListMenusForItem(ctx, itemID)
	if err != nil {
		return nil, fmt.Errorf("list menus for item: %w", err)
	}

	return menus, nil
}

func (f *ItemMenuService) ListItemsInMenu(ctx context.Context, menuID int32) ([]repository.ListItemsInMenuRow, error) {
	items, err := f.queries.ListItemsInMenu(ctx, menuID)
	if err != nil {
		return nil, fmt.Errorf("list items in menu: %w", err)
	}

	return items, nil
}

func (f *ItemMenuService) AddFoodItemToMenu(ctx context.Context, menuID int32, itemID int32) error {
	err := f.queries.AddFoodItemToMenu(ctx, repository.AddFoodItemToMenuParams{
		MenuID: menuID,
		ItemID: itemID,
	})
	if err != nil {
		return fmt.Errorf("add food item to menu: %w", err)
	}

	return nil
}

func (f *ItemMenuService) RemoveFoodItemFromMenu(ctx context.Context, menuID int32, itemID int32) error {
	err := f.queries.RemoveFoodItemFromMenu(ctx, repository.RemoveFoodItemFromMenuParams{
		MenuID: menuID,
		ItemID: itemID,
	})
	if err != nil {
		return fmt.Errorf("remove food item from menu: %w", err)
	}

	return nil
}

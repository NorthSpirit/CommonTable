package foodrelationshipservices

import (
	"context"
	"fmt"

	"CommonTable/internal/repository"
)

type ItemAllergenService struct {
	queries *repository.Queries
}

func NewItemAllergenService(queries *repository.Queries) *ItemAllergenService {
	return &ItemAllergenService{queries: queries}
}

func (f *ItemAllergenService) ListAllergensForItem(ctx context.Context, itemID int32) ([]repository.Allergen, error) {
	allergens, err := f.queries.ListAllergensForItem(ctx, itemID)
	if err != nil {
		return nil, fmt.Errorf("list allergens for item: %w", err)
	}

	return allergens, nil
}

func (f *ItemAllergenService) ListAllergensForItemDisplay(ctx context.Context, itemID int32, locale string) ([]repository.ListAllergensForItemDisplayRow, error) {
	allergens, err := f.queries.ListAllergensForItemDisplay(ctx, repository.ListAllergensForItemDisplayParams{
		Locale: locale,
		ItemID: itemID,
	})
	if err != nil {
		return nil, fmt.Errorf("list allergens for item display: %w", err)
	}

	return allergens, nil
}

func (f *ItemAllergenService) ListItemsWithAllergen(ctx context.Context, allergenID int32) ([]repository.ListItemsWithAllergenRow, error) {
	items, err := f.queries.ListItemsWithAllergen(ctx, allergenID)
	if err != nil {
		return nil, fmt.Errorf("list items with allergen: %w", err)
	}

	return items, nil
}

func (f *ItemAllergenService) AddAllergenToItem(ctx context.Context, itemID int32, allergenID int32) error {
	err := f.queries.AddAllergenToItem(ctx, repository.AddAllergenToItemParams{
		ItemID:     itemID,
		AllergenID: allergenID,
	})
	if err != nil {
		return fmt.Errorf("add allergen to item: %w", err)
	}

	return nil
}

func (f *ItemAllergenService) RemoveAllergenFromItem(ctx context.Context, itemID int32, allergenID int32) error {
	err := f.queries.RemoveAllergenFromItem(ctx, repository.RemoveAllergenFromItemParams{
		ItemID:     itemID,
		AllergenID: allergenID,
	})
	if err != nil {
		return fmt.Errorf("remove allergen from item: %w", err)
	}

	return nil
}

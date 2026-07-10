package foodrelationshipservices

import (
	"CommonTable/internal/helpers"
	"context"
	"fmt"

	"CommonTable/internal/repository"
)

type IngredientItemService struct {
	queries *repository.Queries
}

func NewIngredientItemService(queries *repository.Queries) *IngredientItemService {
	return &IngredientItemService{queries: queries}
}

func (f *IngredientItemService) ListIngredientsForItem(ctx context.Context, itemID int32) ([]repository.ListIngredientsForItemRow, error) {
	ingredients, err := f.queries.ListIngredientsForItem(ctx, itemID)
	if err != nil {
		return nil, fmt.Errorf("list ingredients for item: %w", err)
	}

	return ingredients, nil
}

func (f *IngredientItemService) ListItemsUsingIngredient(ctx context.Context, ingredientID int32) ([]repository.ListItemsUsingIngredientRow, error) {
	items, err := f.queries.ListItemsUsingIngredient(ctx, ingredientID)
	if err != nil {
		return nil, fmt.Errorf("list items using ingredient: %w", err)
	}

	return items, nil
}

func (f *IngredientItemService) AddIngredientToItem(
	ctx context.Context,
	ingredientID int32,
	itemID int32,
	measurement *string,
	amount *float64,
) error {

	err := f.queries.AddIngredientToItem(ctx, repository.AddIngredientToItemParams{
		IngredientID: ingredientID,
		ItemID:       itemID,
		Measurement:  helpers.TextFromString(measurement),
		Amount:       helpers.Float8FromFloat64(amount),
	})
	if err != nil {
		return fmt.Errorf("add ingredient to item: %w", err)
	}

	return nil
}

func (f *IngredientItemService) UpdateIngredientInItem(
	ctx context.Context,
	ingredientID int32,
	itemID int32,
	measurement *string,
	amount *float64,
) error {

	err := f.queries.UpdateIngredientInItem(ctx, repository.UpdateIngredientInItemParams{
		IngredientID: ingredientID,
		ItemID:       itemID,
		Measurement:  helpers.TextFromString(measurement),
		Amount:       helpers.Float8FromFloat64(amount),
	})
	if err != nil {
		return fmt.Errorf("update ingredient in item: %w", err)
	}

	return nil
}

func (f *IngredientItemService) RemoveIngredientFromItem(ctx context.Context, ingredientID int32, itemID int32) error {
	err := f.queries.RemoveIngredientFromItem(ctx, repository.RemoveIngredientFromItemParams{
		IngredientID: ingredientID,
		ItemID:       itemID,
	})
	if err != nil {
		return fmt.Errorf("remove ingredient from item: %w", err)
	}

	return nil
}

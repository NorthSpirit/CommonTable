package foodservices

import (
	"context"
	"fmt"

	"CommonTable/internal/helpers"
	"CommonTable/internal/repository"
)

type FoodIngredientService struct {
	queries *repository.Queries
}

func NewFoodIngredientService(queries *repository.Queries) *FoodIngredientService {
	return &FoodIngredientService{queries: queries}
}

func (f *FoodIngredientService) GetFoodIngredientList(ctx context.Context) ([]repository.FoodIngredient, error) {
	ingredients, err := f.queries.ListFoodIngredients(ctx)
	if err != nil {
		return nil, fmt.Errorf("list food ingredients: %w", err)
	}

	return ingredients, nil
}

func (f *FoodIngredientService) GetFoodIngredientByID(ctx context.Context, id int32) (repository.FoodIngredient, error) {
	ingredient, err := f.queries.GetFoodIngredient(ctx, id)
	if err != nil {
		return repository.FoodIngredient{}, fmt.Errorf("get food ingredient: %w", err)
	}

	return ingredient, nil
}

func (f *FoodIngredientService) InsertFoodIngredient(
	ctx context.Context,
	name string,
	description *string,
) (repository.FoodIngredient, error) {

	ingredient, err := f.queries.CreateFoodIngredient(ctx, repository.CreateFoodIngredientParams{
		Name:        name,
		Description: helpers.TextFromString(description),
	})
	if err != nil {
		return repository.FoodIngredient{}, fmt.Errorf("create food ingredient: %w", err)
	}

	return ingredient, nil
}

func (f *FoodIngredientService) DeleteFoodIngredient(ctx context.Context, id int32) error {
	err := f.queries.DeleteFoodIngredient(ctx, id)
	if err != nil {
		return fmt.Errorf("delete food ingredient: %w", err)
	}

	return nil
}

func (f *FoodIngredientService) UpdateFoodIngredient(
	ctx context.Context,
	id int32,
	name string,
	description *string,
) error {

	err := f.queries.UpdateFoodIngredient(ctx, repository.UpdateFoodIngredientParams{
		ID:          id,
		Name:        name,
		Description: helpers.TextFromString(description),
	})
	if err != nil {
		return fmt.Errorf("update food ingredient: %w", err)
	}

	return nil
}

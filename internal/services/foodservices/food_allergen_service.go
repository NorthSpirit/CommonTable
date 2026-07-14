package foodservices

import (
	"context"
	"fmt"

	"CommonTable/internal/helpers"
	"CommonTable/internal/repository"
)

type FoodAllergenService struct {
	queries *repository.Queries
}

func NewFoodAllergenService(queries *repository.Queries) *FoodAllergenService {
	return &FoodAllergenService{queries: queries}
}

func (f *FoodAllergenService) GetFoodAllergensList(ctx context.Context) ([]repository.Allergen, error) {
	allergens, err := f.queries.ListAllergens(ctx)
	if err != nil {
		return nil, fmt.Errorf("list food allergens: %w", err)
	}

	return allergens, nil
}

func (f *FoodAllergenService) GetFoodAllergenByID(ctx context.Context, id int32) (repository.Allergen, error) {
	allergen, err := f.queries.GetAllergen(ctx, id)
	if err != nil {
		return repository.Allergen{}, fmt.Errorf("get food allergen: %w", err)
	}

	return allergen, nil
}

func (f *FoodAllergenService) GetFoodAllergenByCode(ctx context.Context, code string) (repository.Allergen, error) {
	allergen, err := f.queries.GetAllergenByCode(ctx, code)
	if err != nil {
		return repository.Allergen{}, fmt.Errorf("get food allergen by code: %w", err)
	}

	return allergen, nil
}

func (f *FoodAllergenService) InsertFoodAllergen(
	ctx context.Context,
	name string,
	nameAlt *string,
	description *string,
	descriptionAlt *string,
	code string,
) (repository.Allergen, error) {

	allergen, err := f.queries.CreateAllergen(ctx, repository.CreateAllergenParams{
		Name:           name,
		NameAlt:        helpers.TextFromString(nameAlt),
		Description:    helpers.TextFromString(description),
		DescriptionAlt: helpers.TextFromString(descriptionAlt),
		Code:           code,
	})
	if err != nil {
		return repository.Allergen{}, fmt.Errorf("create food allergen: %w", err)
	}

	return allergen, nil
}

func (f *FoodAllergenService) UpdateFoodAllergen(
	ctx context.Context,
	id int32,
	name string,
	nameAlt *string,
	description *string,
	descriptionAlt *string,
	code string,
) error {

	err := f.queries.UpdateAllergen(ctx, repository.UpdateAllergenParams{
		ID:             id,
		Name:           name,
		NameAlt:        helpers.TextFromString(nameAlt),
		Description:    helpers.TextFromString(description),
		DescriptionAlt: helpers.TextFromString(descriptionAlt),
		Code:           code,
	})
	if err != nil {
		return fmt.Errorf("update food allergen: %w", err)
	}

	return nil
}

func (f *FoodAllergenService) DeleteFoodAllergen(ctx context.Context, id int32) error {
	err := f.queries.DeleteAllergen(ctx, id)
	if err != nil {
		return fmt.Errorf("delete food allergen: %w", err)
	}

	return nil
}

package restaurantservices

import (
	"context"
	"fmt"

	"CommonTable/internal/helpers"
	"CommonTable/internal/repository"
)

type RestaurantService struct {
	queries *repository.Queries
}

func NewRestaurantService(queries *repository.Queries) *RestaurantService {
	return &RestaurantService{queries: queries}
}

func (f *RestaurantService) GetRestaurantByID(ctx context.Context, id int32) (repository.Restaurant, error) {
	restaurant, err := f.queries.GetRestaurant(ctx, id)
	if err != nil {
		return repository.Restaurant{}, fmt.Errorf("get restaurant: %w", err)
	}

	return restaurant, nil
}

// GetRestaurantForDisplay returns locale-resolved restaurant data — for
// customer-facing pages.
func (f *RestaurantService) GetRestaurantForDisplay(ctx context.Context, id int32, locale string) (repository.GetRestaurantForDisplayRow, error) {
	restaurant, err := f.queries.GetRestaurantForDisplay(ctx, repository.GetRestaurantForDisplayParams{
		Locale: locale,
		ID:     id,
	})
	if err != nil {
		return repository.GetRestaurantForDisplayRow{}, fmt.Errorf("get restaurant for display: %w", err)
	}

	return restaurant, nil
}

func (f *RestaurantService) GetRestaurantByName(ctx context.Context, name string) (repository.Restaurant, error) {
	restaurant, err := f.queries.GetRestaurantByName(ctx, name)
	if err != nil {
		return repository.Restaurant{}, fmt.Errorf("get restaurant by name: %w", err)
	}

	return restaurant, nil
}

func (f *RestaurantService) GetRestaurantList(ctx context.Context) ([]repository.Restaurant, error) {
	restaurants, err := f.queries.ListRestaurants(ctx)
	if err != nil {
		return nil, fmt.Errorf("list restaurants: %w", err)
	}

	return restaurants, nil
}

// ListRestaurantsForDisplay returns locale-resolved restaurant data — for
// customer-facing pages.
func (f *RestaurantService) ListRestaurantsForDisplay(ctx context.Context, locale string) ([]repository.ListRestaurantsForDisplayRow, error) {
	restaurants, err := f.queries.ListRestaurantsForDisplay(ctx, locale)
	if err != nil {
		return nil, fmt.Errorf("list restaurants for display: %w", err)
	}

	return restaurants, nil
}

// GetRestaurantListForAdmin includes soft-deleted restaurants, for admin
// screens that need to find and restore them.
func (f *RestaurantService) GetRestaurantListForAdmin(ctx context.Context) ([]repository.Restaurant, error) {
	restaurants, err := f.queries.ListRestaurantsForAdmin(ctx)
	if err != nil {
		return nil, fmt.Errorf("list restaurants for admin: %w", err)
	}

	return restaurants, nil
}

// InsertRestaurant expects passwordHash to already be hashed by the caller —
// this service never handles a raw password.
func (f *RestaurantService) InsertRestaurant(
	ctx context.Context,
	name string,
	address *string,
	phone *string,
	email *string,
	passwordHash string,
	logoURL *string,
	description *string,
	descriptionAlt *string,
) (repository.Restaurant, error) {

	restaurant, err := f.queries.CreateRestaurant(ctx, repository.CreateRestaurantParams{
		Name:           name,
		Address:        helpers.TextFromString(address),
		Phone:          helpers.TextFromString(phone),
		Email:          helpers.TextFromString(email),
		PasswordHash:   helpers.TextFromString(&passwordHash),
		LogoUrl:        helpers.TextFromString(logoURL),
		Description:    helpers.TextFromString(description),
		DescriptionAlt: helpers.TextFromString(descriptionAlt),
	})
	if err != nil {
		return repository.Restaurant{}, fmt.Errorf("create restaurant: %w", err)
	}

	return restaurant, nil
}

func (f *RestaurantService) UpdateRestaurant(
	ctx context.Context,
	id int32,
	name string,
	address *string,
	phone *string,
	email *string,
	logoURL *string,
	description *string,
	descriptionAlt *string,
) error {

	err := f.queries.UpdateRestaurant(ctx, repository.UpdateRestaurantParams{
		ID:             id,
		Name:           name,
		Address:        helpers.TextFromString(address),
		Phone:          helpers.TextFromString(phone),
		Email:          helpers.TextFromString(email),
		LogoUrl:        helpers.TextFromString(logoURL),
		Description:    helpers.TextFromString(description),
		DescriptionAlt: helpers.TextFromString(descriptionAlt),
	})
	if err != nil {
		return fmt.Errorf("update restaurant: %w", err)
	}

	return nil
}

// UpdateRestaurantPasswordHash expects passwordHash to already be hashed
// by the caller.
func (f *RestaurantService) UpdateRestaurantPasswordHash(ctx context.Context, id int32, passwordHash string) error {
	err := f.queries.UpdateRestaurantPasswordHash(ctx, repository.UpdateRestaurantPasswordHashParams{
		ID:           id,
		PasswordHash: helpers.TextFromString(&passwordHash),
	})
	if err != nil {
		return fmt.Errorf("update restaurant password hash: %w", err)
	}

	return nil
}

func (f *RestaurantService) DeleteRestaurant(ctx context.Context, id int32) error {
	err := f.queries.SoftDeleteRestaurant(ctx, id)
	if err != nil {
		return fmt.Errorf("delete restaurant: %w", err)
	}

	return nil
}

func (f *RestaurantService) RestoreRestaurant(ctx context.Context, id int32) error {
	err := f.queries.RestoreRestaurant(ctx, id)
	if err != nil {
		return fmt.Errorf("restore restaurant: %w", err)
	}

	return nil
}

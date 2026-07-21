package customer

import (
	"CommonTable/internal/helpers"
	"CommonTable/internal/repository"
	"CommonTable/internal/templates/shared"
)

func toRestaurantView(r repository.ListRestaurantsForDisplayRow) shared.RestaurantView {
	view := shared.RestaurantView{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
	}
	if r.Address.Valid {
		view.Address = r.Address.String
	}
	if r.LogoUrl.Valid {
		view.LogoURL = r.LogoUrl.String
	}
	return view
}

func toMenuView(m repository.ListFoodMenusByRestaurantForDisplayRow) shared.MenuView {
	view := shared.MenuView{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
	}
	if m.ImageUrl.Valid {
		view.ImageURL = m.ImageUrl.String
	}
	return view
}

func toMenuItemView(item repository.ListMenuItemsForDisplayRow) shared.MenuItemView {
	view := shared.MenuItemView{
		ID:          item.ID,
		Name:        item.Name,
		Description: item.Description,
	}
	if item.ImageUrl.Valid {
		view.ImageURL = item.ImageUrl.String
	}
	if item.PrepTimeMins.Valid {
		view.PrepTimeMins = item.PrepTimeMins.Int32
		view.HasPrepTime = true
	}
	if price, err := helpers.NumericToFloat64(item.Price); err == nil && price != nil {
		view.Price = *price
	}
	return view
}

func toItemDetailView(item repository.GetFoodItemForDisplayRow, allergens []repository.ListAllergensForItemDisplayRow) shared.ItemDetailView {
	view := shared.ItemDetailView{
		ID:           item.ID,
		Name:         item.Name,
		Description:  item.Description,
		Instructions: item.Instructions,
	}
	if item.ImageUrl.Valid {
		view.ImageURL = item.ImageUrl.String
	}
	for _, a := range allergens {
		view.Allergens = append(view.Allergens, a.Name)
	}
	return view
}

func toBrandingView(r repository.Restaurant) shared.BrandingView {
	view := shared.BrandingView{
		Name: r.Name,
	}
	if r.LogoUrl.Valid {
		view.LogoURL = r.LogoUrl.String
	}
	return view
}

func toQuickNavItemView(item repository.ListMostOrderedItemsByRestaurantRow) shared.QuickNavItemView {
	view := shared.QuickNavItemView{
		ID:   item.ID,
		Name: item.Name,
	}
	if item.ImageUrl.Valid {
		view.ImageURL = item.ImageUrl.String
	}
	return view
}

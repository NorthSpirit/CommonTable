package customer

import (
	"CommonTable/internal/helpers"
	"CommonTable/internal/repository"
	"CommonTable/internal/templates"
)

func toRestaurantView(r repository.Restaurant) templates.RestaurantView {
	view := templates.RestaurantView{
		ID:   r.ID,
		Name: r.Name,
	}
	if r.Address.Valid {
		view.Address = r.Address.String
	}
	if r.LogoUrl.Valid {
		view.LogoURL = r.LogoUrl.String
	}
	return view
}

func toMenuView(m repository.FoodMenu) templates.MenuView {
	view := templates.MenuView{
		ID:   m.ID,
		Name: m.Name,
	}
	if m.Description.Valid {
		view.Description = m.Description.String
	}
	return view
}

func toMenuItemView(item repository.ListMenuItemsForDisplayRow) templates.MenuItemView {
	view := templates.MenuItemView{
		ID:   item.ID,
		Name: item.Name,
	}
	if item.Description.Valid {
		view.Description = item.Description.String
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

func toItemDetailView(item repository.GetFoodItemRow, allergens []repository.Allergen) templates.ItemDetailView {
	view := templates.ItemDetailView{
		ID:   item.ID,
		Name: item.Name,
	}
	if item.Description.Valid {
		view.Description = item.Description.String
	}
	if item.Instructions.Valid {
		view.Instructions = item.Instructions.String
	}
	if item.ImageUrl.Valid {
		view.ImageURL = item.ImageUrl.String
	}
	for _, a := range allergens {
		view.Allergens = append(view.Allergens, a.Name)
	}
	return view
}

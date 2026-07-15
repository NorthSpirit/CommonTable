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

func toMenuView(m repository.ListFoodMenusByRestaurantForDisplayRow) templates.MenuView {
	view := templates.MenuView{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
	}
	return view
}

func toMenuItemView(item repository.ListMenuItemsForDisplayRow) templates.MenuItemView {
	view := templates.MenuItemView{
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

func toItemDetailView(item repository.GetFoodItemForDisplayRow, allergens []repository.ListAllergensForItemDisplayRow) templates.ItemDetailView {
	view := templates.ItemDetailView{
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

func toBrandingView(r repository.Restaurant) templates.BrandingView {
	view := templates.BrandingView{
		Name: r.Name,
	}
	if r.LogoUrl.Valid {
		view.LogoURL = r.LogoUrl.String
	}
	return view
}

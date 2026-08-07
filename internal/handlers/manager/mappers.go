package manager

import (
	"CommonTable/internal/helpers"
	"CommonTable/internal/repository"
	managertmpl "CommonTable/internal/templates/manager"
)

func toEmployeeView(e repository.Employee) managertmpl.EmployeeView {
	view := managertmpl.EmployeeView{
		ID:       e.ID,
		Name:     e.Name,
		Role:     string(e.Role),
		IsActive: e.IsActive,
	}
	if e.Email.Valid {
		view.Email = e.Email.String
	}
	if e.PhoneNumber.Valid {
		view.PhoneNumber = e.PhoneNumber.String
	}
	return view
}

func toMenuView(m repository.FoodMenu) managertmpl.MenuView {
	view := managertmpl.MenuView{ID: m.ID, Name: m.Name}
	if m.NameAlt.Valid {
		view.NameAlt = m.NameAlt.String
	}
	if m.Description.Valid {
		view.Description = m.Description.String
	}
	if m.DescriptionAlt.Valid {
		view.DescriptionAlt = m.DescriptionAlt.String
	}
	if m.ImageUrl.Valid {
		view.ImageURL = m.ImageUrl.String
	}
	return view
}

func toFoodItemView(item repository.ListFoodItemsByRestaurantForAdminRow) managertmpl.FoodItemView {
	view := managertmpl.FoodItemView{
		ID:          item.ID,
		Name:        item.Name,
		IsAvailable: item.IsAvailable,
		IsDeleted:   item.DeletedAt.Valid,
	}
	if item.NameAlt.Valid {
		view.NameAlt = item.NameAlt.String
	}
	if item.Description.Valid {
		view.Description = item.Description.String
	}
	if item.DescriptionAlt.Valid {
		view.DescriptionAlt = item.DescriptionAlt.String
	}
	if item.Instructions.Valid {
		view.Instructions = item.Instructions.String
	}
	if item.InstructionsAlt.Valid {
		view.InstructionsAlt = item.InstructionsAlt.String
	}
	if item.ImageUrl.Valid {
		view.ImageURL = item.ImageUrl.String
	}
	if item.PrepTimeMins.Valid {
		view.PrepTimeMins = item.PrepTimeMins.Int32
	}
	if price, err := helpers.NumericToFloat64(item.Price); err == nil && price != nil {
		view.Price = *price
	}
	return view
}

func toFoodItemViewFromGet(item repository.GetFoodItemRow) managertmpl.FoodItemView {
	view := managertmpl.FoodItemView{
		ID:          item.ID,
		Name:        item.Name,
		IsAvailable: item.IsAvailable,
		IsDeleted:   item.DeletedAt.Valid,
	}
	if item.NameAlt.Valid {
		view.NameAlt = item.NameAlt.String
	}
	if item.Description.Valid {
		view.Description = item.Description.String
	}
	if item.DescriptionAlt.Valid {
		view.DescriptionAlt = item.DescriptionAlt.String
	}
	if item.Instructions.Valid {
		view.Instructions = item.Instructions.String
	}
	if item.InstructionsAlt.Valid {
		view.InstructionsAlt = item.InstructionsAlt.String
	}
	if item.ImageUrl.Valid {
		view.ImageURL = item.ImageUrl.String
	}
	if item.PrepTimeMins.Valid {
		view.PrepTimeMins = item.PrepTimeMins.Int32
	}
	if price, err := helpers.NumericToFloat64(item.Price); err == nil && price != nil {
		view.Price = *price
	}
	return view
}

func toFoodItemViewFromMenuAdmin(item repository.ListItemsInMenuForAdminRow) managertmpl.FoodItemView {
	view := managertmpl.FoodItemView{
		ID:          item.ID,
		Name:        item.Name,
		IsAvailable: item.IsAvailable,
		IsDeleted:   item.DeletedAt.Valid,
	}
	if item.ImageUrl.Valid {
		view.ImageURL = item.ImageUrl.String
	}
	if price, err := helpers.NumericToFloat64(item.Price); err == nil && price != nil {
		view.Price = *price
	}
	return view
}

func strPtrOrNil(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func toIngredientOptionView(i repository.FoodIngredient) managertmpl.IngredientOptionView {
	return managertmpl.IngredientOptionView{ID: i.ID, Name: i.Name}
}

func toItemIngredientView(i repository.ListIngredientsForItemRow) managertmpl.ItemIngredientView {
	view := managertmpl.ItemIngredientView{IngredientID: i.ID, Name: i.Name}
	if i.Measurement.Valid {
		view.Measurement = i.Measurement.String
	}
	if i.Amount.Valid {
		view.Amount = i.Amount.Float64
	}
	return view
}

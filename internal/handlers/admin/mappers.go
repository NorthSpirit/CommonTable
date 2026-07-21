package admin

import (
	"CommonTable/internal/helpers"
	"CommonTable/internal/repository"
	"CommonTable/internal/templates/admin/adminshared"
)

func toRestaurantView(r repository.Restaurant) adminshared.RestaurantView {
	view := adminshared.RestaurantView{
		ID:        r.ID,
		Name:      r.Name,
		IsDeleted: r.DeletedAt.Valid,
	}
	if r.Address.Valid {
		view.Address = r.Address.String
	}
	if r.Phone.Valid {
		view.Phone = r.Phone.String
	}
	if r.Email.Valid {
		view.Email = r.Email.String
	}
	if r.LogoUrl.Valid {
		view.LogoURL = r.LogoUrl.String
	}
	if r.Description.Valid {
		view.Description = r.Description.String
	}
	if r.DescriptionAlt.Valid {
		view.DescriptionAlt = r.DescriptionAlt.String
	}
	return view
}

func toFoodItemView(item repository.ListFoodItemsByRestaurantForAdminRow) adminshared.FoodItemView {
	view := adminshared.FoodItemView{
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

func toFoodItemViewFromGet(item repository.GetFoodItemRow) adminshared.FoodItemView {
	view := adminshared.FoodItemView{
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

func toMenuView(m repository.FoodMenu) adminshared.MenuView {
	view := adminshared.MenuView{ID: m.ID, Name: m.Name}
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

func toFoodItemViewFromMenuAdmin(item repository.ListItemsInMenuForAdminRow) adminshared.FoodItemView {
	view := adminshared.FoodItemView{
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

func toIngredientView(i repository.FoodIngredient) adminshared.IngredientView {
	view := adminshared.IngredientView{
		ID:   i.ID,
		Name: i.Name,
	}
	if i.Description.Valid {
		view.Description = i.Description.String
	}
	return view
}

func toAllergenAdminView(a repository.Allergen) adminshared.AllergenAdminView {
	view := adminshared.AllergenAdminView{
		ID:   a.ID,
		Name: a.Name,
		Code: a.Code,
	}
	if a.NameAlt.Valid {
		view.NameAlt = a.NameAlt.String
	}
	if a.Description.Valid {
		view.Description = a.Description.String
	}
	if a.DescriptionAlt.Valid {
		view.DescriptionAlt = a.DescriptionAlt.String
	}
	return view
}

func toEmployeeView(e repository.Employee) adminshared.EmployeeView {
	view := adminshared.EmployeeView{
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

func toTableView(t repository.Table) adminshared.TableView {
	view := adminshared.TableView{
		ID:       t.ID,
		Number:   t.Number,
		Capacity: t.Capacity,
		Status:   string(t.Status),
	}
	if t.Section.Valid {
		view.Section = t.Section.String
	}
	return view
}

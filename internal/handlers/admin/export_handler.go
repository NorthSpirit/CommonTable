package admin

import (
	"encoding/json"
	"net/http"
	"strconv"

	"CommonTable/internal/helpers"
	"CommonTable/internal/services/foodrelationshipservices"
	"CommonTable/internal/services/foodservices"
	"CommonTable/internal/services/restaurantservices"

	"github.com/go-chi/chi/v5"
)

type ExportHandler struct {
	restaurantService     *restaurantservices.RestaurantService
	foodMenuService       *foodservices.FoodMenuService
	itemMenuService       *foodrelationshipservices.ItemMenuService
	itemAllergenService   *foodrelationshipservices.ItemAllergenService
	ingredientItemService *foodrelationshipservices.IngredientItemService
}

func NewExportHandler(
	restaurantService *restaurantservices.RestaurantService,
	foodMenuService *foodservices.FoodMenuService,
	itemMenuService *foodrelationshipservices.ItemMenuService,
	itemAllergenService *foodrelationshipservices.ItemAllergenService,
	ingredientItemService *foodrelationshipservices.IngredientItemService,
) *ExportHandler {
	return &ExportHandler{
		restaurantService: restaurantService, foodMenuService: foodMenuService,
		itemMenuService: itemMenuService, itemAllergenService: itemAllergenService,
		ingredientItemService: ingredientItemService,
	}
}

type exportIngredient struct {
	Name        string  `json:"name"`
	Amount      float64 `json:"amount,omitempty"`
	Measurement string  `json:"measurement,omitempty"`
}

type exportItem struct {
	Name           string             `json:"name"`
	NameAlt        string             `json:"name_alt,omitempty"`
	Description    string             `json:"description,omitempty"`
	DescriptionAlt string             `json:"description_alt,omitempty"`
	Price          float64            `json:"price"`
	ImageURL       string             `json:"image_url,omitempty"`
	IsAvailable    bool               `json:"is_available"`
	Allergens      []string           `json:"allergens"`
	Ingredients    []exportIngredient `json:"ingredients"`
}

type exportMenu struct {
	Name           string       `json:"name"`
	NameAlt        string       `json:"name_alt,omitempty"`
	Description    string       `json:"description,omitempty"`
	DescriptionAlt string       `json:"description_alt,omitempty"`
	ImageURL       string       `json:"image_url,omitempty"`
	Items          []exportItem `json:"items"`
}

type exportRestaurant struct {
	Name           string       `json:"name"`
	Description    string       `json:"description,omitempty"`
	DescriptionAlt string       `json:"description_alt,omitempty"`
	Menus          []exportMenu `json:"menus"`
}

// Export handles GET /admin/restaurants/{id}/export
//
// ID-stripped, self-contained JSON dump of a restaurant's menu content
// — meant for third parties building their own bespoke storefront
// against CommonTable's data, independent of the built-in customer
// web page. See /api-docs/export for field documentation.
func (h *ExportHandler) Export(w http.ResponseWriter, r *http.Request) {
	restaurantID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid restaurant id", http.StatusBadRequest)
		return
	}

	restaurant, err := h.restaurantService.GetRestaurantByID(r.Context(), int32(restaurantID))
	if err != nil {
		http.Error(w, "restaurant not found", http.StatusNotFound)
		return
	}

	menus, err := h.foodMenuService.ListFoodMenusByRestaurant(r.Context(), int32(restaurantID))
	if err != nil {
		http.Error(w, "failed to load menus", http.StatusInternalServerError)
		return
	}

	out := exportRestaurant{Name: restaurant.Name, Menus: []exportMenu{}}
	if restaurant.Description.Valid {
		out.Description = restaurant.Description.String
	}
	if restaurant.DescriptionAlt.Valid {
		out.DescriptionAlt = restaurant.DescriptionAlt.String
	}

	for _, m := range menus {
		menu := exportMenu{Name: m.Name, Items: []exportItem{}}
		if m.NameAlt.Valid {
			menu.NameAlt = m.NameAlt.String
		}
		if m.Description.Valid {
			menu.Description = m.Description.String
		}
		if m.DescriptionAlt.Valid {
			menu.DescriptionAlt = m.DescriptionAlt.String
		}
		if m.ImageUrl.Valid {
			menu.ImageURL = m.ImageUrl.String
		}

		items, err := h.itemMenuService.ListItemsInMenuForAdmin(r.Context(), m.ID)
		if err != nil {
			http.Error(w, "failed to load items", http.StatusInternalServerError)
			return
		}

		for _, item := range items {
			if item.DeletedAt.Valid {
				continue // Never export soft-deleted items.
			}

			exp := exportItem{Name: item.Name, IsAvailable: item.IsAvailable, Allergens: []string{}, Ingredients: []exportIngredient{}}
			if item.ImageUrl.Valid {
				exp.ImageURL = item.ImageUrl.String
			}
			if price, perr := helpers.NumericToFloat64(item.Price); perr == nil && price != nil {
				exp.Price = *price
			}

			allergens, err := h.itemAllergenService.ListAllergensForItem(r.Context(), item.ID)
			if err == nil {
				for _, a := range allergens {
					exp.Allergens = append(exp.Allergens, a.Name)
				}
			}

			ingredients, err := h.ingredientItemService.ListIngredientsForItem(r.Context(), item.ID)
			if err == nil {
				for _, ing := range ingredients {
					expIng := exportIngredient{Name: ing.Name}
					if ing.Amount.Valid {
						expIng.Amount = ing.Amount.Float64
					}
					if ing.Measurement.Valid {
						expIng.Measurement = ing.Measurement.String
					}
					exp.Ingredients = append(exp.Ingredients, expIng)
				}
			}

			menu.Items = append(menu.Items, exp)
		}

		out.Menus = append(out.Menus, menu)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", `attachment; filename="menu-export.json"`)
	if err := json.NewEncoder(w).Encode(out); err != nil {
		http.Error(w, "failed to encode export", http.StatusInternalServerError)
	}
}

package waiterapi

import (
	"CommonTable/internal/services/foodrelationshipservices"
	"net/http"
	"strconv"

	"CommonTable/internal/helpers"
	"CommonTable/internal/services/foodservices"

	"github.com/go-chi/chi/v5"
)

type MenuHandler struct {
	foodMenuService       *foodservices.FoodMenuService
	foodItemService       *foodservices.FoodItemService
	itemAllergenService   *foodrelationshipservices.ItemAllergenService
	ingredientItemService *foodrelationshipservices.IngredientItemService
}

func NewMenuHandler(
	foodMenuService *foodservices.FoodMenuService,
	foodItemService *foodservices.FoodItemService,
	itemAllergenService *foodrelationshipservices.ItemAllergenService,
	ingredientItemService *foodrelationshipservices.IngredientItemService,
) *MenuHandler {
	return &MenuHandler{
		foodMenuService:       foodMenuService,
		foodItemService:       foodItemService,
		itemAllergenService:   itemAllergenService,
		ingredientItemService: ingredientItemService,
	}
}

type menuJSON struct {
	ID          int32  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ListMenus handles GET /api/waiter/menus
func (h *MenuHandler) ListMenus(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromContext(r.Context())

	menus, err := h.foodMenuService.ListFoodMenusByRestaurant(r.Context(), claims.RestaurantID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load menus")
		return
	}

	out := make([]menuJSON, 0, len(menus))
	for _, m := range menus {
		description := ""
		if m.Description.Valid {
			description = m.Description.String
		}
		out = append(out, menuJSON{ID: m.ID, Name: m.Name, Description: description})
	}

	writeJSON(w, http.StatusOK, out)
}

type menuItemJSON struct {
	ID          int32   `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	ImageURL    string  `json:"image_url"`
	IsAvailable bool    `json:"is_available"`
}

// ListMenuItems handles GET /api/waiter/menus/{menuID}/items
func (h *MenuHandler) ListMenuItems(w http.ResponseWriter, r *http.Request) {
	menuID, err := strconv.Atoi(chi.URLParam(r, "menuID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid menu id")
		return
	}

	items, err := h.foodItemService.ListMenuItemsForWaiter(r.Context(), int32(menuID))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load menu items")
		return
	}

	out := make([]menuItemJSON, 0, len(items))
	for _, item := range items {
		description := ""
		if item.Description.Valid {
			description = item.Description.String
		}
		imageURL := ""
		if item.ImageUrl.Valid {
			imageURL = item.ImageUrl.String
		}
		price := 0.0
		if p, err := helpers.NumericToFloat64(item.Price); err == nil && p != nil {
			price = *p
		}
		out = append(out, menuItemJSON{
			ID:          item.ID,
			Name:        item.Name,
			Description: description,
			Price:       price,
			ImageURL:    imageURL,
			IsAvailable: item.IsAvailable,
		})
	}

	writeJSON(w, http.StatusOK, out)
}

type itemIngredientJSON struct {
	Name        string  `json:"name"`
	Measurement string  `json:"measurement"`
	Amount      float64 `json:"amount"`
}

type itemDetailJSON struct {
	ID           int32                `json:"id"`
	Name         string               `json:"name"`
	Description  string               `json:"description"`
	Instructions string               `json:"instructions"`
	ImageURL     string               `json:"image_url"`
	Price        float64              `json:"price"`
	IsAvailable  bool                 `json:"is_available"`
	Ingredients  []itemIngredientJSON `json:"ingredients"`
	Allergens    []string             `json:"allergens"`
}

// ItemDetail handles GET /api/waiter/items/{itemID}
func (h *MenuHandler) ItemDetail(w http.ResponseWriter, r *http.Request) {
	itemID, err := strconv.Atoi(chi.URLParam(r, "itemID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid item id")
		return
	}

	item, err := h.foodItemService.GetFoodItemForDisplay(r.Context(), int32(itemID), "en")
	if err != nil {
		writeError(w, http.StatusNotFound, "item not found")
		return
	}

	ingredients, err := h.ingredientItemService.ListIngredientsForItem(r.Context(), int32(itemID))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load ingredients")
		return
	}

	allergens, err := h.itemAllergenService.ListAllergensForItem(r.Context(), int32(itemID))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load allergens")
		return
	}

	imageURL := ""
	if item.ImageUrl.Valid {
		imageURL = item.ImageUrl.String
	}

	ingredientsOut := make([]itemIngredientJSON, 0, len(ingredients))
	for _, ing := range ingredients {
		measurement := ""
		if ing.Measurement.Valid {
			measurement = ing.Measurement.String
		}
		amount := 0.0
		if ing.Amount.Valid {
			amount = ing.Amount.Float64
		}
		ingredientsOut = append(ingredientsOut, itemIngredientJSON{Name: ing.Name, Measurement: measurement, Amount: amount})
	}

	allergensOut := make([]string, 0, len(allergens))
	for _, a := range allergens {
		allergensOut = append(allergensOut, a.Name)
	}

	// Note: price/availability aren't part of GetFoodItemForDisplay's
	// row — fetched separately so the waiter sees the same price shown
	// on the menu list, without duplicating that lookup logic.
	menuRow, err := h.foodItemService.GetFoodItemByID(r.Context(), int32(itemID))
	price := 0.0
	isAvailable := true
	if err == nil {
		if p, perr := helpers.NumericToFloat64(menuRow.Price); perr == nil && p != nil {
			price = *p
		}
		isAvailable = menuRow.IsAvailable
	}

	writeJSON(w, http.StatusOK, itemDetailJSON{
		ID:           item.ID,
		Name:         item.Name,
		Description:  item.Description,
		Instructions: item.Instructions,
		ImageURL:     imageURL,
		Price:        price,
		IsAvailable:  isAvailable,
		Ingredients:  ingredientsOut,
		Allergens:    allergensOut,
	})
}

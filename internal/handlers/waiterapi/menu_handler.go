package waiterapi

import (
	"net/http"
	"strconv"

	"CommonTable/internal/helpers"
	"CommonTable/internal/services/foodservices"

	"github.com/go-chi/chi/v5"
)

type MenuHandler struct {
	foodMenuService *foodservices.FoodMenuService
	foodItemService *foodservices.FoodItemService
}

func NewMenuHandler(foodMenuService *foodservices.FoodMenuService, foodItemService *foodservices.FoodItemService) *MenuHandler {
	return &MenuHandler{foodMenuService: foodMenuService, foodItemService: foodItemService}
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

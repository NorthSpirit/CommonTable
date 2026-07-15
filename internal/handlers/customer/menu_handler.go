package customer

import (
	"errors"
	"net/http"
	"strconv"

	"CommonTable/internal/services/foodservices"
	"CommonTable/internal/services/restaurantservices"
	"CommonTable/internal/templates"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

type MenuHandler struct {
	foodMenuService   *foodservices.FoodMenuService
	foodItemService   *foodservices.FoodItemService
	restaurantService *restaurantservices.RestaurantService
}

func NewMenuHandler(
	foodMenuService *foodservices.FoodMenuService,
	foodItemService *foodservices.FoodItemService,
	restaurantService *restaurantservices.RestaurantService,
) *MenuHandler {
	return &MenuHandler{
		foodMenuService:   foodMenuService,
		foodItemService:   foodItemService,
		restaurantService: restaurantService,
	}
}

// Items handles GET /menus/{id}
func (h *MenuHandler) Items(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid menu id", http.StatusBadRequest)
		return
	}

	locale := getLocale(r)

	menu, err := h.foodMenuService.GetFoodMenuForDisplay(r.Context(), int32(id), locale)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "failed to load menu", http.StatusInternalServerError)
		return
	}

	restaurant, err := h.restaurantService.GetRestaurantByID(r.Context(), menu.RestaurantID)
	if err != nil {
		http.Error(w, "failed to load restaurant", http.StatusInternalServerError)
		return
	}

	items, err := h.foodItemService.ListMenuItemsForDisplay(r.Context(), int32(id), locale)
	if err != nil {
		http.Error(w, "failed to load menu items", http.StatusInternalServerError)
		return
	}

	views := make([]templates.MenuItemView, 0, len(items))
	for _, item := range items {
		views = append(views, toMenuItemView(item))
	}

	branding := toBrandingView(restaurant)

	if err := templates.MenuItemsPage(branding, menu.Name, views).Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

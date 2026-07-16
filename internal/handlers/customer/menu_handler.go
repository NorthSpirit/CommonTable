package customer

import (
	"errors"
	"net/http"
	"strconv"

	"CommonTable/internal/services/foodservices"
	"CommonTable/internal/services/restaurantservices"
	"CommonTable/internal/templates/client"
	"CommonTable/internal/templates/shared"

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
		renderNotFound(w, r)
		return
	}

	locale := GetLocale(r)

	menu, err := h.foodMenuService.GetFoodMenuForDisplay(r.Context(), int32(id), locale)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			renderNotFound(w, r)
			return
		}
		renderServerError(w, r)
		return
	}

	restaurant, err := h.restaurantService.GetRestaurantByID(r.Context(), menu.RestaurantID)
	if err != nil {
		renderServerError(w, r)
		return
	}

	items, err := h.foodItemService.ListMenuItemsForDisplay(r.Context(), int32(id), locale)
	if err != nil {
		renderServerError(w, r)
		return
	}

	views := make([]shared.MenuItemView, 0, len(items))
	for _, item := range items {
		views = append(views, toMenuItemView(item))
	}

	branding := toBrandingView(restaurant)

	if err := client.MenuItemsPage(locale, branding, menu.Name, views).Render(r.Context(), w); err != nil {
		renderServerError(w, r)
	}
}

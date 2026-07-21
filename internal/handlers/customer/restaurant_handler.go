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

type RestaurantHandler struct {
	restaurantService *restaurantservices.RestaurantService
	foodMenuService   *foodservices.FoodMenuService
	foodItemService   *foodservices.FoodItemService
}

func NewRestaurantHandler(
	restaurantService *restaurantservices.RestaurantService,
	foodMenuService *foodservices.FoodMenuService,
	foodItemService *foodservices.FoodItemService,
) *RestaurantHandler {
	return &RestaurantHandler{
		restaurantService: restaurantService,
		foodMenuService:   foodMenuService,
		foodItemService:   foodItemService,
	}
}

// List handles GET /restaurants
func (h *RestaurantHandler) List(w http.ResponseWriter, r *http.Request) {
	locale := GetLocale(r)

	restaurants, err := h.restaurantService.ListRestaurantsForDisplay(r.Context(), locale)
	if err != nil {
		renderServerError(w, r)
		return
	}

	views := make([]shared.RestaurantView, 0, len(restaurants))
	for _, restaurant := range restaurants {
		views = append(views, toRestaurantView(restaurant))
	}

	if err := client.RestaurantsPage(locale, views).Render(r.Context(), w); err != nil {
		renderServerError(w, r)
	}
}

// Menus handles GET /restaurants/{id}
func (h *RestaurantHandler) Menus(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		renderNotFound(w, r)
		return
	}

	locale := GetLocale(r)

	restaurant, err := h.restaurantService.GetRestaurantByID(r.Context(), int32(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			renderNotFound(w, r)
			return
		}
		renderServerError(w, r)
		return
	}

	menus, err := h.foodMenuService.ListFoodMenusByRestaurantForDisplay(r.Context(), int32(id), locale)
	if err != nil {
		renderServerError(w, r)
		return
	}

	menuViews := make([]shared.MenuView, 0, len(menus))
	for _, menu := range menus {
		menuViews = append(menuViews, toMenuView(menu))
	}

	popularItems, err := h.foodItemService.ListMostOrderedItemsByRestaurant(r.Context(), int32(id), locale, 8)
	if err != nil {
		renderServerError(w, r)
		return
	}

	quickNavViews := make([]shared.QuickNavItemView, 0, len(popularItems))
	for _, item := range popularItems {
		quickNavViews = append(quickNavViews, toQuickNavItemView(item))
	}

	branding := toBrandingView(restaurant)

	if err := client.MenusPage(locale, branding, quickNavViews, menuViews).Render(r.Context(), w); err != nil {
		renderServerError(w, r)
	}
}

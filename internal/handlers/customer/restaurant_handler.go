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

type RestaurantHandler struct {
	restaurantService *restaurantservices.RestaurantService
	foodMenuService   *foodservices.FoodMenuService
}

func NewRestaurantHandler(
	restaurantService *restaurantservices.RestaurantService,
	foodMenuService *foodservices.FoodMenuService,
) *RestaurantHandler {
	return &RestaurantHandler{
		restaurantService: restaurantService,
		foodMenuService:   foodMenuService,
	}
}

// List handles GET /restaurants
func (h *RestaurantHandler) List(w http.ResponseWriter, r *http.Request) {
	restaurants, err := h.restaurantService.GetRestaurantList(r.Context())
	if err != nil {
		http.Error(w, "failed to load restaurants", http.StatusInternalServerError)
		return
	}

	views := make([]templates.RestaurantView, 0, len(restaurants))
	for _, restaurant := range restaurants {
		views = append(views, toRestaurantView(restaurant))
	}

	if err := templates.RestaurantsPage(views).Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

// Menus handles GET /restaurants/{id}
func (h *RestaurantHandler) Menus(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid restaurant id", http.StatusBadRequest)
		return
	}

	locale := getLocale(r)

	restaurant, err := h.restaurantService.GetRestaurantByID(r.Context(), int32(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "failed to load restaurant", http.StatusInternalServerError)
		return
	}

	menus, err := h.foodMenuService.ListFoodMenusByRestaurantForDisplay(r.Context(), int32(id), locale)
	if err != nil {
		http.Error(w, "failed to load menus", http.StatusInternalServerError)
		return
	}

	views := make([]templates.MenuView, 0, len(menus))
	for _, menu := range menus {
		views = append(views, toMenuView(menu))
	}

	branding := toBrandingView(restaurant)

	if err := templates.MenusPage(branding, views).Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

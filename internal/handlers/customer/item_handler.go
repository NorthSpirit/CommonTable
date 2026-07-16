package customer

import (
	"errors"
	"net/http"
	"strconv"

	"CommonTable/internal/services/foodrelationshipservices"
	"CommonTable/internal/services/foodservices"
	"CommonTable/internal/services/restaurantservices"
	"CommonTable/internal/templates/client"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

type ItemHandler struct {
	foodItemService     *foodservices.FoodItemService
	itemAllergenService *foodrelationshipservices.ItemAllergenService
	restaurantService   *restaurantservices.RestaurantService
}

func NewItemHandler(
	foodItemService *foodservices.FoodItemService,
	itemAllergenService *foodrelationshipservices.ItemAllergenService,
	restaurantService *restaurantservices.RestaurantService,
) *ItemHandler {
	return &ItemHandler{
		foodItemService:     foodItemService,
		itemAllergenService: itemAllergenService,
		restaurantService:   restaurantService,
	}
}

// Detail handles GET /items/{id}
func (h *ItemHandler) Detail(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		renderNotFound(w, r)
		return
	}

	locale := GetLocale(r)

	item, err := h.foodItemService.GetFoodItemForDisplay(r.Context(), int32(id), locale)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			renderNotFound(w, r)
			return
		}
		renderServerError(w, r)
		return
	}

	restaurant, err := h.restaurantService.GetRestaurantByID(r.Context(), item.RestaurantID)
	if err != nil {
		renderServerError(w, r)
		return
	}

	allergens, err := h.itemAllergenService.ListAllergensForItemDisplay(r.Context(), int32(id), locale)
	if err != nil {
		renderServerError(w, r)
		return
	}

	view := toItemDetailView(item, allergens)
	branding := toBrandingView(restaurant)

	if err := client.ItemDetailPage(locale, branding, view).Render(r.Context(), w); err != nil {
		renderServerError(w, r)
	}
}

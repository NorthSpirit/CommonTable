package customer

import (
	"errors"
	"net/http"
	"strconv"

	"CommonTable/internal/services/foodrelationshipservices"
	"CommonTable/internal/services/foodservices"
	"CommonTable/internal/services/restaurantservices"
	"CommonTable/internal/templates"

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
		http.Error(w, "invalid item id", http.StatusBadRequest)
		return
	}

	locale := getLocale(r)

	item, err := h.foodItemService.GetFoodItemForDisplay(r.Context(), int32(id), locale)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "failed to load item", http.StatusInternalServerError)
		return
	}

	restaurant, err := h.restaurantService.GetRestaurantByID(r.Context(), item.RestaurantID)
	if err != nil {
		http.Error(w, "failed to load restaurant", http.StatusInternalServerError)
		return
	}

	allergens, err := h.itemAllergenService.ListAllergensForItemDisplay(r.Context(), int32(id), locale)
	if err != nil {
		http.Error(w, "failed to load allergens", http.StatusInternalServerError)
		return
	}

	view := toItemDetailView(item, allergens)
	branding := toBrandingView(restaurant)

	if err := templates.ItemDetailPage(branding, view).Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

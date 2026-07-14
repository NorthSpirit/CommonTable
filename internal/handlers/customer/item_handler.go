package customer

import (
	"errors"
	"net/http"
	"strconv"

	"CommonTable/internal/services/foodrelationshipservices"
	"CommonTable/internal/services/foodservices"
	"CommonTable/internal/templates"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

type ItemHandler struct {
	foodItemService     *foodservices.FoodItemService
	itemAllergenService *foodrelationshipservices.ItemAllergenService
}

func NewItemHandler(
	foodItemService *foodservices.FoodItemService,
	itemAllergenService *foodrelationshipservices.ItemAllergenService,
) *ItemHandler {
	return &ItemHandler{
		foodItemService:     foodItemService,
		itemAllergenService: itemAllergenService,
	}
}

// Detail handles GET /items/{id}
func (h *ItemHandler) Detail(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid item id", http.StatusBadRequest)
		return
	}

	item, err := h.foodItemService.GetFoodItemByID(r.Context(), int32(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "failed to load item", http.StatusInternalServerError)
		return
	}

	allergens, err := h.itemAllergenService.ListAllergensForItem(r.Context(), int32(id))
	if err != nil {
		http.Error(w, "failed to load allergens", http.StatusInternalServerError)
		return
	}

	view := toItemDetailView(item, allergens)

	if err := templates.ItemDetailPage(view).Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

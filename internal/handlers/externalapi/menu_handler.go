package externalapi

import (
	"net/http"

	"CommonTable/internal/helpers"
	"CommonTable/internal/services/externalorderservices"
)

type MenuHandler struct {
	orders *externalorderservices.ExternalOrderService
}

func NewMenuHandler(orders *externalorderservices.ExternalOrderService) *MenuHandler {
	return &MenuHandler{orders: orders}
}

type itemJSON struct {
	ID          int32   `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	ImageURL    string  `json:"image_url"`
}

// ListItems handles GET /api/external/items
func (h *MenuHandler) ListItems(w http.ResponseWriter, r *http.Request) {
	key := apiKeyFromContext(r.Context())

	items, err := h.orders.ListAvailableItems(r.Context(), key.RestaurantID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load items")
		return
	}

	out := make([]itemJSON, 0, len(items))
	for _, item := range items {
		price := 0.0
		if p, err := helpers.NumericToFloat64(item.Price); err == nil && p != nil {
			price = *p
		}
		out = append(out, itemJSON{
			ID:          item.ID,
			Name:        item.Name,
			Description: item.Description.String,
			Price:       price,
			ImageURL:    item.ImageUrl.String,
		})
	}

	writeJSON(w, http.StatusOK, out)
}

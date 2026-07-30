package waiterapi

import (
	"CommonTable/internal/helpers"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"CommonTable/internal/services/waiterservices"

	"github.com/go-chi/chi/v5"
)

type QuickItemsHandler struct {
	quickItemService *waiterservices.WaiterQuickItemService
}

func NewQuickItemsHandler(quickItemService *waiterservices.WaiterQuickItemService) *QuickItemsHandler {
	return &QuickItemsHandler{quickItemService: quickItemService}
}

type quickItemJSON struct {
	ItemID    int32   `json:"item_id"`
	Name      string  `json:"name"`
	ImageURL  string  `json:"image_url"`
	Price     float64 `json:"price"`
	SortOrder int32   `json:"sort_order"`
}

// List handles GET /api/waiter/quick-items
func (h *QuickItemsHandler) List(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromContext(r.Context())

	items, err := h.quickItemService.ListQuickItems(r.Context(), claims.EmployeeID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load quick items")
		return
	}

	out := make([]quickItemJSON, 0, len(items))
	for _, item := range items {
		imageURL := ""
		if item.ImageUrl.Valid {
			imageURL = item.ImageUrl.String
		}
		price := 0.0
		if p, err := helpers.NumericToFloat64(item.Price); err == nil && p != nil {
			price = *p
		}
		out = append(out, quickItemJSON{ItemID: item.ID, Name: item.Name, ImageURL: imageURL, Price: price, SortOrder: item.SortOrder})
	}

	writeJSON(w, http.StatusOK, out)
}

type addQuickItemRequest struct {
	ItemID    int32 `json:"item_id"`
	SortOrder int32 `json:"sort_order"`
}

// Add handles POST /api/waiter/quick-items
func (h *QuickItemsHandler) Add(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromContext(r.Context())

	var req addQuickItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.ItemID == 0 {
		writeError(w, http.StatusBadRequest, "item_id is required")
		return
	}

	err := h.quickItemService.AddQuickItem(r.Context(), claims.EmployeeID, req.ItemID, req.SortOrder)
	if err != nil {
		if errors.Is(err, waiterservices.ErrQuickItemLimitReached) {
			writeError(w, http.StatusConflict, "quick item limit reached (max 15)")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to add quick item")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]bool{"success": true})
}

// Remove handles DELETE /api/waiter/quick-items/{itemID}
func (h *QuickItemsHandler) Remove(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromContext(r.Context())

	itemID, err := strconv.Atoi(chi.URLParam(r, "itemID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid item id")
		return
	}

	if err := h.quickItemService.RemoveQuickItem(r.Context(), claims.EmployeeID, int32(itemID)); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to remove quick item")
		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{"success": true})
}

type reorderRequest struct {
	SortOrder int32 `json:"sort_order"`
}

// Reorder handles PATCH /api/waiter/quick-items/{itemID}
func (h *QuickItemsHandler) Reorder(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromContext(r.Context())

	itemID, err := strconv.Atoi(chi.URLParam(r, "itemID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid item id")
		return
	}

	var req reorderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.quickItemService.UpdateSortOrder(r.Context(), claims.EmployeeID, int32(itemID), req.SortOrder); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update sort order")
		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{"success": true})
}

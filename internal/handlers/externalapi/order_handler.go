package externalapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"CommonTable/internal/helpers"
	"CommonTable/internal/services/externalorderservices"
)

type OrderHandler struct {
	orders *externalorderservices.ExternalOrderService
}

func NewOrderHandler(orders *externalorderservices.ExternalOrderService) *OrderHandler {
	return &OrderHandler{orders: orders}
}

type orderLineRequest struct {
	ItemID   int32  `json:"item_id"`
	Quantity int32  `json:"quantity"`
	Notes    string `json:"notes"`
}

type createOrderRequest struct {
	IdempotencyKey string             `json:"idempotency_key"`
	CustomerName   string             `json:"customer_name"`
	CustomerPhone  string             `json:"customer_phone"`
	Notes          string             `json:"notes"`
	Items          []orderLineRequest `json:"items"`
}

type createOrderResponse struct {
	OrderID   int32  `json:"order_id"`
	Status    string `json:"status"`
	Duplicate bool   `json:"duplicate"`
}

const maxLines = 50

// phoneAllowed accepts digits and the usual separators.
func phoneAllowed(phone string) bool {
	for _, r := range phone {
		if (r < '0' || r > '9') && !strings.ContainsRune("+-() ", r) {
			return false
		}
	}
	return true
}

// Create handles POST /api/external/orders
func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	key := apiKeyFromContext(r.Context())

	var req createOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	req.CustomerName = strings.TrimSpace(req.CustomerName)
	req.CustomerPhone = strings.TrimSpace(req.CustomerPhone)

	if req.IdempotencyKey == "" {
		writeError(w, http.StatusBadRequest, "idempotency_key is required")
		return
	}
	if req.CustomerName == "" || req.CustomerPhone == "" {
		writeError(w, http.StatusBadRequest, "customer_name and customer_phone are required")
		return
	}
	if len(req.Items) == 0 || len(req.Items) > maxLines {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("items must contain 1 to %d lines", maxLines))
		return
	}
	if !phoneAllowed(req.CustomerPhone) {
		writeError(w, http.StatusBadRequest, "customer_phone may contain only digits, +, -, ( ) and spaces")
		return
	}

	if msg := helpers.CheckLengths(
		helpers.FieldLimit{Label: "idempotency_key", Value: req.IdempotencyKey, Max: 64},
		helpers.FieldLimit{Label: "customer_name", Value: req.CustomerName, Max: 255},
		helpers.FieldLimit{Label: "customer_phone", Value: req.CustomerPhone, Max: 30},
		helpers.FieldLimit{Label: "notes", Value: req.Notes, Max: 500},
	); msg != "" {
		writeError(w, http.StatusBadRequest, msg)
		return
	}

	lines := make([]externalorderservices.OrderLineInput, 0, len(req.Items))
	for _, item := range req.Items {
		if msg := helpers.CheckRange("quantity", float64(item.Quantity), 1, 1000); msg != "" {
			writeError(w, http.StatusBadRequest, msg)
			return
		}
		if msg := helpers.CheckLengths(helpers.FieldLimit{Label: "item notes", Value: item.Notes, Max: 500}); msg != "" {
			writeError(w, http.StatusBadRequest, msg)
			return
		}
		lines = append(lines, externalorderservices.OrderLineInput{ItemID: item.ItemID, Quantity: item.Quantity, Notes: item.Notes})
	}

	result, err := h.orders.Submit(r.Context(), key, externalorderservices.SubmitInput{
		IdempotencyKey: req.IdempotencyKey,
		CustomerName:   req.CustomerName,
		CustomerPhone:  req.CustomerPhone,
		Notes:          req.Notes,
		Lines:          lines,
	})
	if err != nil {
		var itemErr *externalorderservices.ItemError
		switch {
		case errors.Is(err, externalorderservices.ErrRateLimited):
			w.Header().Set("Retry-After", "60")
			writeError(w, http.StatusTooManyRequests, "too many orders in a short time, try again later")
		case errors.Is(err, externalorderservices.ErrTooManyPending),
			errors.Is(err, externalorderservices.ErrPhoneLimit):
			writeError(w, http.StatusTooManyRequests, "too many orders are waiting for approval, try again later")
		case errors.As(err, &itemErr) && errors.Is(err, externalorderservices.ErrItemNotFound):
			writeError(w, http.StatusBadRequest, fmt.Sprintf("item %d not found", itemErr.ItemID))
		case errors.As(err, &itemErr) && errors.Is(err, externalorderservices.ErrItemUnavailable):
			writeError(w, http.StatusConflict, fmt.Sprintf("item %d is not available", itemErr.ItemID))
		default:
			writeError(w, http.StatusInternalServerError, "failed to submit order")
		}
		return
	}

	status := http.StatusCreated
	if result.Duplicate {
		status = http.StatusOK
	}
	writeJSON(w, status, createOrderResponse{OrderID: result.OrderID, Status: "pending", Duplicate: result.Duplicate})
}

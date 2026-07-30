package waiterapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"CommonTable/internal/repository"
	"CommonTable/internal/services/foodorderservices"
	"CommonTable/internal/services/foodservices"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderHandler struct {
	pool             *pgxpool.Pool
	foodOrderService *foodorderservices.FoodOrderService
	itemOrderService *foodorderservices.ItemOrderService
	foodItemService  *foodservices.FoodItemService
}

func NewOrderHandler(
	pool *pgxpool.Pool,
	foodOrderService *foodorderservices.FoodOrderService,
	itemOrderService *foodorderservices.ItemOrderService,
	foodItemService *foodservices.FoodItemService,
) *OrderHandler {
	return &OrderHandler{
		pool:             pool,
		foodOrderService: foodOrderService,
		itemOrderService: itemOrderService,
		foodItemService:  foodItemService,
	}
}

type orderLineRequest struct {
	ItemID   int32  `json:"item_id"`
	Quantity int32  `json:"quantity"`
	Notes    string `json:"notes"`
}

type createOrderRequest struct {
	IdempotencyKey string             `json:"idempotency_key"`
	TableID        int32              `json:"table_id"`
	Notes          string             `json:"notes"`
	Items          []orderLineRequest `json:"items"`
}

type createOrderResponse struct {
	OrderID     int32  `json:"order_id"`
	OrderNumber int32  `json:"order_number"`
	Status      string `json:"status"`
	Duplicate   bool   `json:"duplicate"`
	Appended    bool   `json:"appended"`
}

// Create handles POST /api/waiter/orders
//
// Behavior:
//   - If the table already has an open order (not delivered/cancelled),
//     the submitted items are appended to it instead of creating a new
//     order — this is what turns "2 more beers" into an addition to the
//     existing tab rather than a second bill.
//   - Otherwise, a new order is created.
//   - Both paths are idempotency-key protected: a retried request with
//     the same key is safe and won't duplicate items or orders.
func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromContext(r.Context())

	var req createOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.IdempotencyKey == "" {
		writeError(w, http.StatusBadRequest, "idempotency_key is required")
		return
	}
	if req.TableID == 0 {
		writeError(w, http.StatusBadRequest, "table_id is required")
		return
	}
	if len(req.Items) == 0 {
		writeError(w, http.StatusBadRequest, "at least one item is required")
		return
	}

	lines := make([]foodorderservices.OrderLine, 0, len(req.Items))
	for _, item := range req.Items {
		lines = append(lines, foodorderservices.OrderLine{
			ItemID:   item.ItemID,
			Quantity: item.Quantity,
			Notes:    strPtrOrNilJSON(item.Notes),
		})
	}

	// Check whether this table already has an open order — if so,
	// append to it instead of creating a duplicate.
	activeOrder, activeErr := h.foodOrderService.GetActiveOrderForTable(r.Context(), req.TableID, claims.RestaurantID)
	if activeErr == nil {
		if err := h.foodOrderService.AppendItemsToOrderAtomic(r.Context(), h.pool, activeOrder.ID, req.IdempotencyKey, lines); err != nil {
			writeError(w, http.StatusInternalServerError, "failed to add items to existing order")
			return
		}

		writeJSON(w, http.StatusOK, createOrderResponse{
			OrderID:     activeOrder.ID,
			OrderNumber: activeOrder.OrderNumber,
			Status:      string(activeOrder.OrderStatus),
			Duplicate:   false,
			Appended:    true,
		})
		return
	}
	if !errors.Is(activeErr, pgx.ErrNoRows) {
		writeError(w, http.StatusInternalServerError, "failed to check for active order")
		return
	}

	// No active order for this table — check for a prior successful
	// submission with this exact key (a retry), then create fresh.
	existing, err := h.foodOrderService.GetOrderByIdempotencyKey(r.Context(), req.IdempotencyKey)
	if err == nil {
		writeJSON(w, http.StatusOK, createOrderResponse{
			OrderID:     existing.ID,
			OrderNumber: existing.OrderNumber,
			Status:      string(existing.OrderStatus),
			Duplicate:   true,
		})
		return
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		writeError(w, http.StatusInternalServerError, "failed to check for existing order")
		return
	}

	orderNumber, err := h.foodOrderService.GetNextOrderNumberForToday(r.Context(), claims.RestaurantID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to generate order number")
		return
	}

	notes := strPtrOrNilJSON(req.Notes)

	order, err := h.foodOrderService.CreateOrderWithItemsAtomic(
		r.Context(), h.pool, claims.RestaurantID, req.TableID, claims.EmployeeID, orderNumber, notes, req.IdempotencyKey, lines,
	)
	if err != nil {
		if isUniqueViolation(err) {
			existing, existErr := h.foodOrderService.GetOrderByIdempotencyKey(r.Context(), req.IdempotencyKey)
			if existErr == nil {
				writeJSON(w, http.StatusOK, createOrderResponse{
					OrderID:     existing.ID,
					OrderNumber: existing.OrderNumber,
					Status:      string(existing.OrderStatus),
					Duplicate:   true,
				})
				return
			}
		}
		writeError(w, http.StatusInternalServerError, "failed to create order")
		return
	}

	writeJSON(w, http.StatusCreated, createOrderResponse{
		OrderID:     order.ID,
		OrderNumber: order.OrderNumber,
		Status:      string(order.OrderStatus),
		Duplicate:   false,
		Appended:    false,
	})
}

// RemoveItem handles DELETE /api/waiter/orders/{orderID}/items/{itemOrderID}
//
// Only permitted while the order is still "pending" — once the kitchen
// has started preparing it, removing a line would mean un-cooking food,
// which isn't something this endpoint can or should pretend to solve.
func (h *OrderHandler) RemoveItem(w http.ResponseWriter, r *http.Request) {
	orderID, itemOrderID, ok := parseOrderAndItemOrderID(r)
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid order or item id")
		return
	}

	order, err := h.foodOrderService.GetFoodOrder(r.Context(), orderID)
	if err != nil {
		writeError(w, http.StatusNotFound, "order not found")
		return
	}
	if order.OrderStatus != repository.OrderStatusPending {
		writeError(w, http.StatusConflict, "order is already being prepared and can no longer be edited")
		return
	}

	if err := h.itemOrderService.RemoveItemFromOrder(r.Context(), itemOrderID); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to remove item")
		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{"success": true})
}

type updateQuantityRequest struct {
	Quantity int32 `json:"quantity"`
}

// UpdateItemQuantity handles PATCH /api/waiter/orders/{orderID}/items/{itemOrderID}
// Same "pending only" gate as RemoveItem.
func (h *OrderHandler) UpdateItemQuantity(w http.ResponseWriter, r *http.Request) {
	orderID, itemOrderID, ok := parseOrderAndItemOrderID(r)
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid order or item id")
		return
	}

	var req updateQuantityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Quantity <= 0 {
		writeError(w, http.StatusBadRequest, "quantity must be greater than 0")
		return
	}

	order, err := h.foodOrderService.GetFoodOrder(r.Context(), orderID)
	if err != nil {
		writeError(w, http.StatusNotFound, "order not found")
		return
	}
	if order.OrderStatus != repository.OrderStatusPending {
		writeError(w, http.StatusConflict, "order is already being prepared and can no longer be edited")
		return
	}

	if err := h.itemOrderService.UpdateItemOrderQuantity(r.Context(), itemOrderID, req.Quantity); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update quantity")
		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{"success": true})
}

func parseOrderAndItemOrderID(r *http.Request) (orderID int32, itemOrderID int32, ok bool) {
	oid, err1 := strconv.Atoi(chi.URLParam(r, "orderID"))
	iid, err2 := strconv.Atoi(chi.URLParam(r, "itemOrderID"))
	if err1 != nil || err2 != nil {
		return 0, 0, false
	}
	return int32(oid), int32(iid), true
}

func strPtrOrNilJSON(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

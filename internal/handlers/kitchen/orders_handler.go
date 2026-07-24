package kitchen

import (
	"encoding/json"
	"net/http"
	"strconv"

	"CommonTable/internal/helpers"
	"CommonTable/internal/repository"
	"CommonTable/internal/services/foodorderservices"
	"CommonTable/internal/templates/kitchen"
)

type OrdersHandler struct {
	foodOrderService *foodorderservices.FoodOrderService
	itemOrderService *foodorderservices.ItemOrderService
}

func NewOrdersHandler(foodOrderService *foodorderservices.FoodOrderService, itemOrderService *foodorderservices.ItemOrderService) *OrdersHandler {
	return &OrdersHandler{foodOrderService: foodOrderService, itemOrderService: itemOrderService}
}

func nextStatus(current repository.OrderStatus) (repository.OrderStatus, string) {
	switch current {
	case repository.OrderStatusPending:
		return repository.OrderStatusPreparing, "Start Preparing"
	case repository.OrderStatusPreparing:
		return repository.OrderStatusReady, "Mark Ready"
	case repository.OrderStatusReady:
		return repository.OrderStatusDelivered, "Mark Delivered"
	default:
		return "", ""
	}
}

func (h *OrdersHandler) buildBoard(r *http.Request) ([]kitchen.KitchenOrderView, error) {
	claims := claimsFromContext(r.Context())

	orders, err := h.foodOrderService.ListActiveOrdersByRestaurantForKitchen(r.Context(), claims.RestaurantID)
	if err != nil {
		return nil, err
	}

	views := make([]kitchen.KitchenOrderView, 0, len(orders))
	for _, o := range orders {
		items, err := h.itemOrderService.ListItemsForOrder(r.Context(), o.ID)
		if err != nil {
			return nil, err
		}

		itemViews := make([]kitchen.KitchenOrderItemView, 0, len(items))
		for _, item := range items {
			name := "Item"
			if item.ItemNameSnapshot.Valid {
				name = item.ItemNameSnapshot.String
			}
			notes := ""
			if item.Notes.Valid {
				notes = item.Notes.String
			}
			itemViews = append(itemViews, kitchen.KitchenOrderItemView{
				Name:     name,
				Quantity: item.Quantity,
				Notes:    notes,
			})
		}

		next, label := nextStatus(o.OrderStatus)

		notes := ""
		if o.Notes.Valid {
			notes = o.Notes.String
		}

		views = append(views, kitchen.KitchenOrderView{
			ID:              o.ID,
			OrderNumber:     o.OrderNumber,
			TableNumber:     o.TableNumber,
			Status:          string(o.OrderStatus),
			NextStatus:      string(next),
			NextStatusLabel: label,
			CreatedAt:       helpers.FormatSentAt(o.OrderCreationDate),
			Notes:           notes,
			Items:           itemViews,
		})
	}

	return views, nil
}

// Board handles GET /kitchen/orders
func (h *OrdersHandler) Board(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromContext(r.Context())

	views, err := h.buildBoard(r)
	if err != nil {
		http.Error(w, "failed to load orders", http.StatusInternalServerError)
		return
	}

	if err := kitchen.BoardPage(claims.Name, views).Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

type kitchenItemJSON struct {
	Name     string `json:"name"`
	Quantity int32  `json:"quantity"`
	Notes    string `json:"notes"`
}

type kitchenOrderJSON struct {
	ID              int32             `json:"id"`
	OrderNumber     int32             `json:"order_number"`
	TableNumber     int32             `json:"table_number"`
	Status          string            `json:"status"`
	NextStatus      string            `json:"next_status"`
	NextStatusLabel string            `json:"next_status_label"`
	CreatedAt       string            `json:"created_at"`
	Notes           string            `json:"notes"`
	Items           []kitchenItemJSON `json:"items"`
}

// Poll handles GET /kitchen/orders/poll
func (h *OrdersHandler) Poll(w http.ResponseWriter, r *http.Request) {
	views, err := h.buildBoard(r)
	if err != nil {
		http.Error(w, "failed to load orders", http.StatusInternalServerError)
		return
	}

	out := make([]kitchenOrderJSON, 0, len(views))
	for _, v := range views {
		items := make([]kitchenItemJSON, 0, len(v.Items))
		for _, item := range v.Items {
			items = append(items, kitchenItemJSON{Name: item.Name, Quantity: item.Quantity, Notes: item.Notes})
		}
		out = append(out, kitchenOrderJSON{
			ID:              v.ID,
			OrderNumber:     v.OrderNumber,
			TableNumber:     v.TableNumber,
			Status:          v.Status,
			NextStatus:      v.NextStatus,
			NextStatusLabel: v.NextStatusLabel,
			CreatedAt:       v.CreatedAt,
			Notes:           v.Notes,
			Items:           items,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(out)
}

// UpdateStatus handles POST /kitchen/orders/status
func (h *OrdersHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form submission", http.StatusBadRequest)
		return
	}

	orderID, err := strconv.Atoi(r.FormValue("order_id"))
	if err != nil {
		http.Error(w, "invalid order id", http.StatusBadRequest)
		return
	}

	status := repository.OrderStatus(r.FormValue("status"))

	if err := h.foodOrderService.UpdateOrderStatus(r.Context(), int32(orderID), status); err != nil {
		http.Error(w, "failed to update order status", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/kitchen/orders", http.StatusFound)
}

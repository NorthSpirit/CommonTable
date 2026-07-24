package admin

import (
	"errors"
	"net/http"
	"strconv"

	"CommonTable/internal/helpers"
	"CommonTable/internal/repository"
	"CommonTable/internal/services/foodorderservices"
	"CommonTable/internal/templates/admin/adminshared"
	"CommonTable/internal/templates/admin/order"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

type OrderHandler struct {
	foodOrderService *foodorderservices.FoodOrderService
	itemOrderService *foodorderservices.ItemOrderService
	paymentService   *foodorderservices.PaymentService
}

func NewOrderHandler(
	foodOrderService *foodorderservices.FoodOrderService,
	itemOrderService *foodorderservices.ItemOrderService,
	paymentService *foodorderservices.PaymentService,
) *OrderHandler {
	return &OrderHandler{foodOrderService: foodOrderService, itemOrderService: itemOrderService, paymentService: paymentService}
}

// List handles GET /admin/restaurants/{id}/orders
func (h *OrderHandler) List(w http.ResponseWriter, r *http.Request) {
	restaurantID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid restaurant id", http.StatusBadRequest)
		return
	}

	orders, err := h.foodOrderService.ListOrdersByRestaurantForAdmin(r.Context(), int32(restaurantID))
	if err != nil {
		http.Error(w, "failed to load orders", http.StatusInternalServerError)
		return
	}

	views := make([]adminshared.OrderListItemView, 0, len(orders))
	for _, o := range orders {
		views = append(views, adminshared.OrderListItemView{
			ID:          o.ID,
			OrderNumber: o.OrderNumber,
			TableNumber: o.TableNumber,
			WaiterName:  o.WaiterName,
			Status:      string(o.OrderStatus),
			CreatedAt:   helpers.FormatSentAt(o.OrderCreationDate),
		})
	}

	if err := order.OrdersPage(int32(restaurantID), "", views).Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

func (h *OrderHandler) buildDetailView(r *http.Request, orderID int32) (adminshared.OrderDetailView, error) {
	o, err := h.foodOrderService.GetOrderForAdmin(r.Context(), orderID)
	if err != nil {
		return adminshared.OrderDetailView{}, err
	}

	items, err := h.itemOrderService.ListItemsForOrder(r.Context(), orderID)
	if err != nil {
		return adminshared.OrderDetailView{}, err
	}

	itemViews := make([]adminshared.OrderDetailItemView, 0, len(items))
	for _, item := range items {
		name := "Item"
		if item.ItemNameSnapshot.Valid {
			name = item.ItemNameSnapshot.String
		}
		notes := ""
		if item.Notes.Valid {
			notes = item.Notes.String
		}
		price := 0.0
		if p, err := helpers.NumericToFloat64(item.UnitPriceSnapshot); err == nil && p != nil {
			price = *p
		}
		itemViews = append(itemViews, adminshared.OrderDetailItemView{Name: name, Quantity: item.Quantity, Price: price, Notes: notes})
	}

	payments, err := h.paymentService.ListPaymentsByOrder(r.Context(), orderID)
	if err != nil {
		return adminshared.OrderDetailView{}, err
	}

	paymentViews := make([]adminshared.OrderPaymentView, 0, len(payments))
	for _, p := range payments {
		amount := 0.0
		if a, err := helpers.NumericToFloat64(p.Amount); err == nil && a != nil {
			amount = *a
		}
		paymentViews = append(paymentViews, adminshared.OrderPaymentView{ID: p.ID, Amount: amount, Method: string(p.Method), Status: string(p.Status)})
	}

	notes := ""
	if o.Notes.Valid {
		notes = o.Notes.String
	}

	return adminshared.OrderDetailView{
		ID:          o.ID,
		OrderNumber: o.OrderNumber,
		TableNumber: o.TableNumber,
		WaiterName:  o.WaiterName,
		Status:      string(o.OrderStatus),
		CreatedAt:   helpers.FormatSentAt(o.OrderCreationDate),
		Notes:       notes,
		Items:       itemViews,
		Payments:    paymentViews,
	}, nil
}

// Detail handles GET /admin/restaurants/{id}/orders/{orderID}
func (h *OrderHandler) Detail(w http.ResponseWriter, r *http.Request) {
	restaurantID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid restaurant id", http.StatusBadRequest)
		return
	}
	orderID, err := strconv.Atoi(chi.URLParam(r, "orderID"))
	if err != nil {
		http.Error(w, "invalid order id", http.StatusBadRequest)
		return
	}

	view, err := h.buildDetailView(r, int32(orderID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "failed to load order", http.StatusInternalServerError)
		return
	}

	if err := order.OrderDetailPage(int32(restaurantID), view, "").Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

// Update handles POST /admin/restaurants/{id}/orders/{orderID}
func (h *OrderHandler) Update(w http.ResponseWriter, r *http.Request) {
	restaurantID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid restaurant id", http.StatusBadRequest)
		return
	}
	orderID, err := strconv.Atoi(chi.URLParam(r, "orderID"))
	if err != nil {
		http.Error(w, "invalid order id", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form submission", http.StatusBadRequest)
		return
	}

	status := repository.OrderStatus(r.FormValue("status"))
	if err := h.foodOrderService.UpdateOrderStatus(r.Context(), int32(orderID), status); err != nil {
		http.Error(w, "failed to update order status", http.StatusInternalServerError)
		return
	}

	notes := r.FormValue("notes")
	if err := h.foodOrderService.UpdateFoodOrderNotes(r.Context(), int32(orderID), strPtrOrNil(notes)); err != nil {
		http.Error(w, "failed to update order notes", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/restaurants/"+strconv.Itoa(restaurantID)+"/orders/"+strconv.Itoa(orderID), http.StatusFound)
}

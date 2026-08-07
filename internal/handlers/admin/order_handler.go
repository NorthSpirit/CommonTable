package admin

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"CommonTable/internal/helpers"
	"CommonTable/internal/repository"
	"CommonTable/internal/services/foodorderservices"
	"CommonTable/internal/services/foodrelationshipservices"
	"CommonTable/internal/services/foodservices"
	"CommonTable/internal/services/restaurantservices"
	"CommonTable/internal/templates/admin/adminshared"
	"CommonTable/internal/templates/admin/order"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderHandler struct {
	pool              *pgxpool.Pool
	foodOrderService  *foodorderservices.FoodOrderService
	itemOrderService  *foodorderservices.ItemOrderService
	paymentService    *foodorderservices.PaymentService
	tableService      *foodorderservices.TableService
	employeeService   *restaurantservices.EmployeeService
	restaurantService *restaurantservices.RestaurantService
	foodMenuService   *foodservices.FoodMenuService
	itemMenuService   *foodrelationshipservices.ItemMenuService
}

func NewOrderHandler(
	pool *pgxpool.Pool,
	foodOrderService *foodorderservices.FoodOrderService,
	itemOrderService *foodorderservices.ItemOrderService,
	paymentService *foodorderservices.PaymentService,
	tableService *foodorderservices.TableService,
	employeeService *restaurantservices.EmployeeService,
	restaurantService *restaurantservices.RestaurantService,
	foodMenuService *foodservices.FoodMenuService,
	itemMenuService *foodrelationshipservices.ItemMenuService,
) *OrderHandler {
	return &OrderHandler{
		pool: pool, foodOrderService: foodOrderService, itemOrderService: itemOrderService,
		paymentService: paymentService, tableService: tableService, employeeService: employeeService,
		restaurantService: restaurantService, foodMenuService: foodMenuService, itemMenuService: itemMenuService,
	}
}

// List handles GET /admin/restaurants/{id}/orders
func (h *OrderHandler) List(w http.ResponseWriter, r *http.Request) {
	restaurantID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid restaurant id", http.StatusBadRequest)
		return
	}

	restaurant, err := h.restaurantService.GetRestaurantByID(r.Context(), int32(restaurantID))
	if err != nil {
		http.Error(w, "failed to load restaurant", http.StatusInternalServerError)
		return
	}

	orders, err := h.foodOrderService.ListOrdersByRestaurantForAdmin(r.Context(), int32(restaurantID))
	if err != nil {
		http.Error(w, "failed to load orders", http.StatusInternalServerError)
		return
	}

	views := make([]adminshared.OrderListItemView, 0, len(orders))
	for _, o := range orders {
		tableNumber := int32(0)
		if o.TableNumber.Valid {
			tableNumber = o.TableNumber.Int32
		}

		views = append(views, adminshared.OrderListItemView{
			ID:          o.ID,
			OrderNumber: o.OrderNumber,
			TableNumber: tableNumber,
			WaiterName:  o.WaiterName,
			Status:      string(o.OrderStatus),
			CreatedAt:   helpers.FormatSentAt(o.OrderCreationDate),
		})
	}

	if err := order.OrdersPage(int32(restaurantID), restaurant.Name, views).Render(r.Context(), w); err != nil {
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
		itemViews = append(itemViews, adminshared.OrderDetailItemView{ID: item.ID, Name: name, Quantity: item.Quantity, Price: price, Notes: notes})
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

	tableNumber := int32(0)
	if o.TableNumber.Valid {
		tableNumber = o.TableNumber.Int32
	}

	total := 0.0
	for _, item := range itemViews {
		total += item.Price * float64(item.Quantity)
	}

	totalPaid := 0.0
	for _, p := range paymentViews {
		if p.Status == "paid" {
			totalPaid += p.Amount
		}
	}

	remaining := total - totalPaid
	if remaining < 0 {
		remaining = 0
	}

	return adminshared.OrderDetailView{
		ID:          o.ID,
		OrderNumber: o.OrderNumber,
		TableNumber: tableNumber,
		WaiterName:  o.WaiterName,
		Status:      string(o.OrderStatus),
		CreatedAt:   helpers.FormatSentAt(o.OrderCreationDate),
		Notes:       notes,
		Items:       itemViews,
		Payments:    paymentViews,
		Total:       total,
		TotalPaid:   totalPaid,
		Remaining:   remaining,
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

	// Cancelling an order that already has a paid payment would leave
	// money taken against an order that says it never happened — block
	// this and require the payment to be explicitly marked refunded
	// first, rather than silently moving money on a status change.
	if status == repository.OrderStatusCancelled {
		payments, err := h.paymentService.ListPaymentsByOrder(r.Context(), int32(orderID))
		if err == nil {
			for _, p := range payments {
				if p.Status == repository.PaymentStatusPaid {
					view, _ := h.buildDetailView(r, int32(orderID))
					w.WriteHeader(http.StatusConflict)
					_ = order.OrderDetailPage(int32(restaurantID), view, "Cannot cancel: this order has a paid payment. Mark it refunded first, then cancel.").Render(r.Context(), w)
					return
				}
			}
		}
	}

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

// UpdatePaymentStatus handles POST /admin/restaurants/{id}/orders/{orderID}/payments/{paymentID}
func (h *OrderHandler) UpdatePaymentStatus(w http.ResponseWriter, r *http.Request) {
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
	paymentID, err := strconv.Atoi(chi.URLParam(r, "paymentID"))
	if err != nil {
		http.Error(w, "invalid payment id", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form submission", http.StatusBadRequest)
		return
	}

	status := repository.PaymentStatus(r.FormValue("status"))
	if err := h.paymentService.UpdatePaymentStatus(r.Context(), int32(paymentID), status); err != nil {
		http.Error(w, "failed to update payment status", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/restaurants/"+strconv.Itoa(restaurantID)+"/orders/"+strconv.Itoa(orderID), http.StatusFound)
}

func (h *OrderHandler) loadFormOptions(r *http.Request, restaurantID int32) ([]adminshared.TableOptionView, []adminshared.EmployeeOptionView, []adminshared.MenuItemOptionView, error) {
	tables, err := h.tableService.ListTables(r.Context())
	if err != nil {
		return nil, nil, nil, err
	}
	tableOptions := make([]adminshared.TableOptionView, 0, len(tables))
	for _, t := range tables {
		tableOptions = append(tableOptions, adminshared.TableOptionView{ID: t.ID, Number: t.Number})
	}

	employees, err := h.employeeService.ListActiveEmployeesByRestaurant(r.Context(), restaurantID)
	if err != nil {
		return nil, nil, nil, err
	}
	waiterOptions := make([]adminshared.EmployeeOptionView, 0)
	for _, e := range employees {
		if string(e.Role) == "waiter" || string(e.Role) == "manager" {
			waiterOptions = append(waiterOptions, adminshared.EmployeeOptionView{ID: e.ID, Name: e.Name})
		}
	}

	menus, err := h.foodMenuService.ListFoodMenusByRestaurant(r.Context(), restaurantID)
	if err != nil {
		return nil, nil, nil, err
	}
	itemOptions := make([]adminshared.MenuItemOptionView, 0)
	seen := map[int32]bool{}
	for _, m := range menus {
		items, err := h.itemMenuService.ListItemsInMenuForAdmin(r.Context(), m.ID)
		if err != nil {
			continue
		}
		for _, item := range items {
			if seen[item.ID] || item.DeletedAt.Valid || !item.IsAvailable {
				continue
			}
			seen[item.ID] = true
			price := 0.0
			if p, err := helpers.NumericToFloat64(item.Price); err == nil && p != nil {
				price = *p
			}
			itemOptions = append(itemOptions, adminshared.MenuItemOptionView{ID: item.ID, Name: item.Name, Price: price})
		}
	}

	return tableOptions, waiterOptions, itemOptions, nil
}

// NewForm handles GET /admin/restaurants/{id}/orders/new
func (h *OrderHandler) NewForm(w http.ResponseWriter, r *http.Request) {
	restaurantID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid restaurant id", http.StatusBadRequest)
		return
	}

	tableOptions, waiterOptions, itemOptions, err := h.loadFormOptions(r, int32(restaurantID))
	if err != nil {
		http.Error(w, "failed to load form options", http.StatusInternalServerError)
		return
	}

	if err := order.OrderFormPage(int32(restaurantID), tableOptions, waiterOptions, itemOptions, "").Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

// Create handles POST /admin/restaurants/{id}/orders/new
//
// table_id of "0" from the form means Take-Out (no table) — see
// order_form_page.templ's dropdown. A takeout order is never merged
// into an existing one, since multiple unrelated customers could be
// ordering takeout at the same time.
func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	restaurantID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid restaurant id", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form submission", http.StatusBadRequest)
		return
	}

	renderFormError := func(message string) {
		tableOptions, waiterOptions, itemOptions, loadErr := h.loadFormOptions(r, int32(restaurantID))
		if loadErr != nil {
			http.Error(w, "failed to load form options", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		_ = order.OrderFormPage(int32(restaurantID), tableOptions, waiterOptions, itemOptions, message).Render(r.Context(), w)
	}

	tableIDRaw, err := strconv.Atoi(r.FormValue("table_id"))
	if err != nil {
		renderFormError("Invalid table selection.")
		return
	}

	var tableIDPtr *int32
	if tableIDRaw != 0 {
		v := int32(tableIDRaw)
		tableIDPtr = &v
	}

	waiterID, err := strconv.Atoi(r.FormValue("waiter_id"))
	if err != nil {
		renderFormError("Invalid waiter selection.")
		return
	}

	var lines []foodorderservices.OrderLine
	for key, values := range r.PostForm {
		if !strings.HasPrefix(key, "quantity_") || len(values) == 0 {
			continue
		}
		itemIDStr := strings.TrimPrefix(key, "quantity_")
		itemID, err := strconv.Atoi(itemIDStr)
		if err != nil {
			continue
		}
		qty, err := strconv.Atoi(values[0])
		if err != nil || qty <= 0 {
			continue
		}
		lines = append(lines, foodorderservices.OrderLine{ItemID: int32(itemID), Quantity: int32(qty)})
	}

	if len(lines) == 0 {
		renderFormError("At least one item with quantity above 0 is required.")
		return
	}

	idempotencyKey := fmt.Sprintf("admin-%d-%d", restaurantID, time.Now().UnixNano())

	if tableIDPtr != nil {
		activeOrder, activeErr := h.foodOrderService.GetActiveOrderForTable(r.Context(), *tableIDPtr, int32(restaurantID))
		if activeErr == nil {
			if err := h.foodOrderService.AppendItemsToOrderAtomic(r.Context(), h.pool, activeOrder.ID, idempotencyKey, lines); err != nil {
				http.Error(w, "failed to add items to existing order", http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/admin/restaurants/"+strconv.Itoa(restaurantID)+"/orders", http.StatusFound)
			return
		}
		if !errors.Is(activeErr, pgx.ErrNoRows) {
			http.Error(w, "failed to check for active order", http.StatusInternalServerError)
			return
		}
	}

	notes := strPtrOrNil(r.FormValue("notes"))

	orderNumber, err := h.foodOrderService.GetNextOrderNumberForToday(r.Context(), int32(restaurantID))
	if err != nil {
		http.Error(w, "failed to generate order number", http.StatusInternalServerError)
		return
	}

	_, err = h.foodOrderService.CreateOrderWithItemsAtomic(
		r.Context(), h.pool, int32(restaurantID), tableIDPtr, int32(waiterID), orderNumber, notes, idempotencyKey, lines,
	)
	if err != nil {
		http.Error(w, "failed to create order", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/restaurants/"+strconv.Itoa(restaurantID)+"/orders", http.StatusFound)
}

// RemoveItem handles POST /admin/restaurants/{id}/orders/{orderID}/items/{itemOrderID}/delete
//
// Same gate as the waiter app: only permitted while the order is still
// "pending" — once the kitchen has started preparing it, removing a
// line would mean un-cooking food.
func (h *OrderHandler) RemoveItem(w http.ResponseWriter, r *http.Request) {
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
	itemOrderID, err := strconv.Atoi(chi.URLParam(r, "itemOrderID"))
	if err != nil {
		http.Error(w, "invalid item id", http.StatusBadRequest)
		return
	}

	order, err := h.foodOrderService.GetFoodOrder(r.Context(), int32(orderID))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	if order.OrderStatus != repository.OrderStatusPending {
		http.Error(w, "order is already being prepared and can no longer be edited", http.StatusConflict)
		return
	}

	if err := h.itemOrderService.RemoveItemFromOrder(r.Context(), int32(itemOrderID)); err != nil {
		http.Error(w, "failed to remove item", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/restaurants/"+strconv.Itoa(restaurantID)+"/orders/"+strconv.Itoa(orderID), http.StatusFound)
}

// UpdateItemQuantity handles POST /admin/restaurants/{id}/orders/{orderID}/items/{itemOrderID}/edit
// Same "pending only" gate as RemoveItem.
func (h *OrderHandler) UpdateItemQuantity(w http.ResponseWriter, r *http.Request) {
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
	itemOrderID, err := strconv.Atoi(chi.URLParam(r, "itemOrderID"))
	if err != nil {
		http.Error(w, "invalid item id", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form submission", http.StatusBadRequest)
		return
	}

	quantity, err := strconv.Atoi(r.FormValue("quantity"))
	if err != nil || quantity <= 0 {
		http.Error(w, "quantity must be greater than 0", http.StatusBadRequest)
		return
	}

	order, err := h.foodOrderService.GetFoodOrder(r.Context(), int32(orderID))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	if order.OrderStatus != repository.OrderStatusPending {
		http.Error(w, "order is already being prepared and can no longer be edited", http.StatusConflict)
		return
	}

	if err := h.itemOrderService.UpdateItemOrderQuantity(r.Context(), int32(itemOrderID), int32(quantity)); err != nil {
		http.Error(w, "failed to update quantity", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/restaurants/"+strconv.Itoa(restaurantID)+"/orders/"+strconv.Itoa(orderID), http.StatusFound)
}

// CreatePayment handles POST /admin/restaurants/{id}/orders/{orderID}/payments/new
func (h *OrderHandler) CreatePayment(w http.ResponseWriter, r *http.Request) {
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

	amount, err := strconv.ParseFloat(r.FormValue("amount"), 64)
	if err != nil || amount <= 0 {
		http.Error(w, "a valid amount is required", http.StatusBadRequest)
		return
	}

	method := repository.PaymentMethod(r.FormValue("method"))
	status := repository.PaymentStatus(r.FormValue("status"))

	now := time.Now()
	if _, err := h.paymentService.InsertPayment(r.Context(), int32(orderID), amount, method, nil, &now, status); err != nil {
		http.Error(w, "failed to record payment", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/restaurants/"+strconv.Itoa(restaurantID)+"/orders/"+strconv.Itoa(orderID), http.StatusFound)
}

// QuickPayCash handles POST /admin/restaurants/{id}/orders/{orderID}/payments/quick-pay
//
// One-click "mark fully paid in cash" — computes the remaining balance
// server-side (never trusts a client-submitted amount for this action)
// and records it as a single paid cash payment.
func (h *OrderHandler) QuickPayCash(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "failed to load order", http.StatusInternalServerError)
		return
	}

	if view.Remaining <= 0 {
		http.Redirect(w, r, "/admin/restaurants/"+strconv.Itoa(restaurantID)+"/orders/"+strconv.Itoa(orderID), http.StatusFound)
		return
	}

	now := time.Now()
	if _, err := h.paymentService.InsertPayment(r.Context(), int32(orderID), view.Remaining, repository.PaymentMethodCash, nil, &now, repository.PaymentStatusPaid); err != nil {
		http.Error(w, "failed to record payment", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/restaurants/"+strconv.Itoa(restaurantID)+"/orders/"+strconv.Itoa(orderID), http.StatusFound)
}

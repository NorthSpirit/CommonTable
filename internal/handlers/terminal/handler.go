package terminal

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"CommonTable/internal/helpers"
	"CommonTable/internal/repository"
	"CommonTable/internal/services/authservices"
	"CommonTable/internal/services/foodorderservices"
	terminaltmpl "CommonTable/internal/templates/terminal"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

type Handler struct {
	authService      *authservices.AuthService
	foodOrderService *foodorderservices.FoodOrderService
	itemOrderService *foodorderservices.ItemOrderService
	paymentService   *foodorderservices.PaymentService
}

func NewHandler(
	authService *authservices.AuthService,
	foodOrderService *foodorderservices.FoodOrderService,
	itemOrderService *foodorderservices.ItemOrderService,
	paymentService *foodorderservices.PaymentService,
) *Handler {
	return &Handler{authService: authService, foodOrderService: foodOrderService, itemOrderService: itemOrderService, paymentService: paymentService}
}

// LoginPage handles GET /terminal/login
func (h *Handler) LoginPage(w http.ResponseWriter, r *http.Request) {
	if err := terminaltmpl.LoginPage("").Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

// Login handles POST /terminal/login
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form submission", http.StatusBadRequest)
		return
	}

	name := r.FormValue("restaurant_name")
	password := r.FormValue("password")

	restaurant, err := h.authService.VerifyRestaurantPassword(r.Context(), name, password)
	if err != nil {
		if errors.Is(err, authservices.ErrInvalidCredentials) || errors.Is(err, pgx.ErrNoRows) {
			w.WriteHeader(http.StatusBadRequest)
			_ = terminaltmpl.LoginPage("Incorrect restaurant name or password.").Render(r.Context(), w)
			return
		}
		http.Error(w, "login failed", http.StatusInternalServerError)
		return
	}

	setSessionCookie(w, strconv.Itoa(int(restaurant.ID)))
	http.Redirect(w, r, "/terminal", http.StatusFound)
}

// Logout handles GET /terminal/logout
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	clearSessionCookie(w)
	http.Redirect(w, r, "/terminal/login", http.StatusFound)
}

// Board handles GET /terminal
func (h *Handler) Board(w http.ResponseWriter, r *http.Request) {
	restaurantID := restaurantIDFromContext(r.Context())

	orders, err := h.foodOrderService.ListOrdersByRestaurantForAdmin(r.Context(), restaurantID)
	if err != nil {
		http.Error(w, "failed to load orders", http.StatusInternalServerError)
		return
	}

	views := make([]terminaltmpl.PendingOrderView, 0)
	for _, o := range orders {
		if o.OrderStatus == repository.OrderStatusCancelled {
			continue
		}

		items, err := h.itemOrderService.ListItemsForOrder(r.Context(), o.ID)
		if err != nil {
			continue
		}
		total := 0.0
		for _, item := range items {
			if p, err := helpers.NumericToFloat64(item.UnitPriceSnapshot); err == nil && p != nil {
				total += *p * float64(item.Quantity)
			}
		}

		payments, err := h.paymentService.ListPaymentsByOrder(r.Context(), o.ID)
		if err != nil {
			continue
		}
		paid := 0.0
		for _, p := range payments {
			if p.Status == repository.PaymentStatusPaid {
				if amt, err := helpers.NumericToFloat64(p.Amount); err == nil && amt != nil {
					paid += *amt
				}
			}
		}

		remaining := total - paid
		if remaining <= 0 {
			continue
		}

		tableLabel := "Take-Out"
		if o.TableNumber.Valid {
			tableLabel = "Table " + strconv.Itoa(int(o.TableNumber.Int32))
		}

		views = append(views, terminaltmpl.PendingOrderView{
			OrderID:     o.ID,
			OrderNumber: o.OrderNumber,
			TableLabel:  tableLabel,
			Remaining:   remaining,
		})
	}

	restaurantName := "Restaurant"

	if err := terminaltmpl.BoardPage(restaurantName, views).Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

// Simulate handles POST /terminal/orders/{orderID}/simulate
//
// Stands in for what a real terminal would do: contact the processor
// directly, then the processor calls our server back (a webhook) with
// the result. Here, the "tap" and the "confirmation" happen in one
// request since there's no real processor involved — but the shape
// (create pending, then resolve to paid/failed) mirrors the real flow.
func (h *Handler) Simulate(w http.ResponseWriter, r *http.Request) {
	orderID, err := strconv.Atoi(chi.URLParam(r, "orderID"))
	if err != nil {
		http.Error(w, "invalid order id", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form submission", http.StatusBadRequest)
		return
	}

	items, err := h.itemOrderService.ListItemsForOrder(r.Context(), int32(orderID))
	if err != nil {
		http.Error(w, "failed to load order items", http.StatusInternalServerError)
		return
	}
	total := 0.0
	for _, item := range items {
		if p, err := helpers.NumericToFloat64(item.UnitPriceSnapshot); err == nil && p != nil {
			total += *p * float64(item.Quantity)
		}
	}

	payments, _ := h.paymentService.ListPaymentsByOrder(r.Context(), int32(orderID))
	paid := 0.0
	for _, p := range payments {
		if p.Status == repository.PaymentStatusPaid {
			if amt, err := helpers.NumericToFloat64(p.Amount); err == nil && amt != nil {
				paid += *amt
			}
		}
	}
	remaining := total - paid

	status := repository.PaymentStatusFailed
	if r.FormValue("outcome") == "success" {
		status = repository.PaymentStatusPaid
	}

	now := time.Now()
	if _, err := h.paymentService.InsertPayment(r.Context(), int32(orderID), remaining, repository.PaymentMethodCard, nil, &now, status); err != nil {
		http.Error(w, "failed to record payment", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/terminal", http.StatusFound)
}

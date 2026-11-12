package manager

import (
	"net/http"
	"strconv"
	"strings"

	"CommonTable/internal/helpers"
	"CommonTable/internal/services/externalorderservices"
	"CommonTable/internal/services/restaurantservices"
	managertmpl "CommonTable/internal/templates/manager"

	"github.com/go-chi/chi/v5"
)

var onlineOrderMessages = map[string]string{
	"approved": "Order approved. It is now in the Orders list.",
	"rejected": "Order rejected.",
	"stale":    "That order was already decided.",
}

type OnlineOrderHandler struct {
	orders            *externalorderservices.ExternalOrderService
	restaurantService *restaurantservices.RestaurantService
}

func NewOnlineOrderHandler(orders *externalorderservices.ExternalOrderService, restaurantService *restaurantservices.RestaurantService) *OnlineOrderHandler {
	return &OnlineOrderHandler{orders: orders, restaurantService: restaurantService}
}

// List handles GET /manager/orders/online
func (h *OnlineOrderHandler) List(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromContext(r.Context())

	restaurant, err := h.restaurantService.GetRestaurantByID(r.Context(), claims.RestaurantID)
	if err != nil {
		http.Error(w, "failed to load restaurant", http.StatusInternalServerError)
		return
	}

	pending, err := h.orders.ListPendingWithItems(r.Context(), claims.RestaurantID)
	if err != nil {
		http.Error(w, "failed to load online orders", http.StatusInternalServerError)
		return
	}

	views := make([]managertmpl.OnlineOrderView, 0, len(pending))
	for _, p := range pending {
		items := make([]managertmpl.OnlineOrderItemView, 0, len(p.Items))
		for _, item := range p.Items {
			price := 0.0
			if v, err := helpers.NumericToFloat64(item.UnitPriceSnapshot); err == nil && v != nil {
				price = *v
			}
			items = append(items, managertmpl.OnlineOrderItemView{
				Name:      item.ItemNameSnapshot,
				Quantity:  item.Quantity,
				UnitPrice: strconv.FormatFloat(price, 'f', 2, 64),
				Notes:     item.Notes.String,
			})
		}

		views = append(views, managertmpl.OnlineOrderView{
			ID:            p.Order.ID,
			CustomerName:  p.Order.CustomerName,
			CustomerPhone: p.Order.CustomerPhone,
			Notes:         p.Order.Notes.String,
			CreatedAt:     helpers.FormatSentAt(p.Order.CreatedAt),
			Total:         strconv.FormatFloat(p.Total, 'f', 2, 64),
			Items:         items,
		})
	}

	message := onlineOrderMessages[r.URL.Query().Get("done")]

	if err := managertmpl.OnlineOrdersPage(restaurant.Name, claims.Name, message, views).Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

// Approve handles POST /manager/orders/online/{onlineOrderID}/approve
func (h *OnlineOrderHandler) Approve(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromContext(r.Context())

	id, err := strconv.Atoi(chi.URLParam(r, "onlineOrderID"))
	if err != nil {
		http.Error(w, "invalid order id", http.StatusBadRequest)
		return
	}

	if _, err := h.orders.Approve(r.Context(), claims.RestaurantID, int32(id), claims.EmployeeID); err != nil {
		if err == externalorderservices.ErrAlreadyDecided {
			http.Redirect(w, r, "/manager/orders/online?done=stale", http.StatusSeeOther)
			return
		}
		http.Error(w, "failed to approve order", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/manager/orders/online?done=approved", http.StatusSeeOther)
}

// Reject handles POST /manager/orders/online/{onlineOrderID}/reject
func (h *OnlineOrderHandler) Reject(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromContext(r.Context())

	id, err := strconv.Atoi(chi.URLParam(r, "onlineOrderID"))
	if err != nil {
		http.Error(w, "invalid order id", http.StatusBadRequest)
		return
	}

	reason := strings.TrimSpace(r.FormValue("reason"))
	if msg := helpers.CheckLengths(helpers.FieldLimit{Label: "Reason", Value: reason, Max: 500}); msg != "" {
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	if err := h.orders.Reject(r.Context(), claims.RestaurantID, int32(id), claims.EmployeeID, reason); err != nil {
		if err == externalorderservices.ErrAlreadyDecided {
			http.Redirect(w, r, "/manager/orders/online?done=stale", http.StatusSeeOther)
			return
		}
		http.Error(w, "failed to reject order", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/manager/orders/online?done=rejected", http.StatusSeeOther)
}

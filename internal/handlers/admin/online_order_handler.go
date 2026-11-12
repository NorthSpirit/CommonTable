package admin

import (
	"net/http"
	"strconv"
	"strings"

	"CommonTable/internal/helpers"
	"CommonTable/internal/services/externalorderservices"
	"CommonTable/internal/services/restaurantservices"
	"CommonTable/internal/templates/admin/order"

	"github.com/go-chi/chi/v5"
)

var onlineOrderMessages = map[string]string{
	"approved": "Order approved. It is now in the Orders list.",
	"rejected": "Order rejected.",
	"stale":    "That order was already decided.",
}

type OnlineOrderHandler struct {
	orders            *externalorderservices.ExternalOrderService
	employeeService   *restaurantservices.EmployeeService
	restaurantService *restaurantservices.RestaurantService
}

func NewOnlineOrderHandler(
	orders *externalorderservices.ExternalOrderService,
	employeeService *restaurantservices.EmployeeService,
	restaurantService *restaurantservices.RestaurantService,
) *OnlineOrderHandler {
	return &OnlineOrderHandler{orders: orders, employeeService: employeeService, restaurantService: restaurantService}
}

func onlineListURL(restaurantID int, done string) string {
	return "/admin/restaurants/" + strconv.Itoa(restaurantID) + "/orders/online?done=" + done
}

// List handles GET /admin/restaurants/{id}/orders/online
func (h *OnlineOrderHandler) List(w http.ResponseWriter, r *http.Request) {
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

	pending, err := h.orders.ListPendingWithItems(r.Context(), int32(restaurantID))
	if err != nil {
		http.Error(w, "failed to load online orders", http.StatusInternalServerError)
		return
	}

	employees, err := h.employeeService.ListActiveEmployeesByRestaurant(r.Context(), int32(restaurantID))
	if err != nil {
		http.Error(w, "failed to load employees", http.StatusInternalServerError)
		return
	}
	managers := make([]order.ManagerOptionView, 0)
	for _, e := range employees {
		if string(e.Role) == "manager" {
			managers = append(managers, order.ManagerOptionView{ID: e.ID, Name: e.Name})
		}
	}

	views := make([]order.OnlineOrderView, 0, len(pending))
	for _, p := range pending {
		items := make([]order.OnlineOrderItemView, 0, len(p.Items))
		for _, item := range p.Items {
			price := 0.0
			if v, err := helpers.NumericToFloat64(item.UnitPriceSnapshot); err == nil && v != nil {
				price = *v
			}
			items = append(items, order.OnlineOrderItemView{
				Name:      item.ItemNameSnapshot,
				Quantity:  item.Quantity,
				UnitPrice: strconv.FormatFloat(price, 'f', 2, 64),
				Notes:     item.Notes.String,
			})
		}

		views = append(views, order.OnlineOrderView{
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

	if err := order.OnlineOrdersPage(int32(restaurantID), restaurant.Name, message, managers, views).Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

// Approve handles POST /admin/restaurants/{id}/orders/online/{onlineOrderID}/approve
// The admin picks which of the restaurant's managers the order is approved as.
func (h *OnlineOrderHandler) Approve(w http.ResponseWriter, r *http.Request) {
	restaurantID, err1 := strconv.Atoi(chi.URLParam(r, "id"))
	onlineOrderID, err2 := strconv.Atoi(chi.URLParam(r, "onlineOrderID"))
	if err1 != nil || err2 != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	managerID, err := strconv.Atoi(r.FormValue("employee_id"))
	if err != nil {
		http.Error(w, "Select a manager to approve as.", http.StatusBadRequest)
		return
	}

	employees, err := h.employeeService.ListActiveEmployeesByRestaurant(r.Context(), int32(restaurantID))
	if err != nil {
		http.Error(w, "failed to load employees", http.StatusInternalServerError)
		return
	}
	valid := false
	for _, e := range employees {
		if int(e.ID) == managerID && string(e.Role) == "manager" {
			valid = true
			break
		}
	}
	if !valid {
		http.Error(w, "That is not an active manager of this restaurant.", http.StatusBadRequest)
		return
	}

	if _, err := h.orders.Approve(r.Context(), int32(restaurantID), int32(onlineOrderID), int32(managerID)); err != nil {
		if err == externalorderservices.ErrAlreadyDecided {
			http.Redirect(w, r, onlineListURL(restaurantID, "stale"), http.StatusSeeOther)
			return
		}
		http.Error(w, "failed to approve order", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, onlineListURL(restaurantID, "approved"), http.StatusSeeOther)
}

// Reject handles POST /admin/restaurants/{id}/orders/online/{onlineOrderID}/reject
func (h *OnlineOrderHandler) Reject(w http.ResponseWriter, r *http.Request) {
	restaurantID, err1 := strconv.Atoi(chi.URLParam(r, "id"))
	onlineOrderID, err2 := strconv.Atoi(chi.URLParam(r, "onlineOrderID"))
	if err1 != nil || err2 != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	reason := strings.TrimSpace(r.FormValue("reason"))
	if msg := helpers.CheckLengths(helpers.FieldLimit{Label: "Reason", Value: reason, Max: 500}); msg != "" {
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	if err := h.orders.Reject(r.Context(), int32(restaurantID), int32(onlineOrderID), 0, reason); err != nil {
		if err == externalorderservices.ErrAlreadyDecided {
			http.Redirect(w, r, onlineListURL(restaurantID, "stale"), http.StatusSeeOther)
			return
		}
		http.Error(w, "failed to reject order", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, onlineListURL(restaurantID, "rejected"), http.StatusSeeOther)
}

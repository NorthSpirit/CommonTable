package admin

import (
	"net/http"
	"strconv"

	"CommonTable/internal/services/foodorderservices"
	"CommonTable/internal/services/restaurantservices"
	"CommonTable/internal/templates/admin/adminshared"
	"CommonTable/internal/templates/admin/order"

	"github.com/go-chi/chi/v5"
)

type TableStatusHandler struct {
	tableService      *foodorderservices.TableService
	restaurantService *restaurantservices.RestaurantService
}

func NewTableStatusHandler(tableService *foodorderservices.TableService, restaurantService *restaurantservices.RestaurantService) *TableStatusHandler {
	return &TableStatusHandler{tableService: tableService, restaurantService: restaurantService}
}

// Page handles GET /admin/restaurants/{id}/table-status
func (h *TableStatusHandler) Page(w http.ResponseWriter, r *http.Request) {
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

	tables, err := h.tableService.ListTablesWithOrderStatusForRestaurant(r.Context(), int32(restaurantID))
	if err != nil {
		http.Error(w, "failed to load tables", http.StatusInternalServerError)
		return
	}

	views := make([]adminshared.TableStatusView, 0, len(tables))
	for _, t := range tables {
		views = append(views, adminshared.TableStatusView{
			ID:       t.ID,
			Number:   t.Number,
			Capacity: t.Capacity,
			Occupied: t.OccupiedByRestaurant,
		})
	}

	if err := order.TableStatusPage(int32(restaurantID), restaurant.Name, views).Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

package manager

import (
	"net/http"

	"CommonTable/internal/services/restaurantservices"
	"CommonTable/internal/templates/manager"
)

type DashboardHandler struct {
	restaurantService *restaurantservices.RestaurantService
}

func NewDashboardHandler(restaurantService *restaurantservices.RestaurantService) *DashboardHandler {
	return &DashboardHandler{restaurantService: restaurantService}
}

// Dashboard handles GET /manager
func (h *DashboardHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromContext(r.Context())

	restaurant, err := h.restaurantService.GetRestaurantByID(r.Context(), claims.RestaurantID)
	if err != nil {
		http.Error(w, "failed to load restaurant", http.StatusInternalServerError)
		return
	}

	if err := manager.DashboardPage(restaurant.Name, claims.Name).Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

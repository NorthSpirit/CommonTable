package waiterapi

import "net/http"

type MeHandler struct{}

func NewMeHandler() *MeHandler {
	return &MeHandler{}
}

type meResponse struct {
	EmployeeID   int32  `json:"employee_id"`
	Name         string `json:"name"`
	Role         string `json:"role"`
	RestaurantID int32  `json:"restaurant_id"`
}

// Me handles GET /api/waiter/me
func (h *MeHandler) Me(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromContext(r.Context())
	writeJSON(w, http.StatusOK, meResponse{
		EmployeeID:   claims.EmployeeID,
		Name:         claims.Name,
		Role:         claims.Role,
		RestaurantID: claims.RestaurantID,
	})
}

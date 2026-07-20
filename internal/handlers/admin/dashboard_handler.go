package admin

import (
	"net/http"

	"CommonTable/internal/templates/admin"
)

type DashboardHandler struct{}

func NewDashboardHandler() *DashboardHandler {
	return &DashboardHandler{}
}

func (h *DashboardHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	if err := admin.DashboardPage().Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

package hub

import (
	"net/http"

	"CommonTable/internal/templates/hub"
)

type HubHandler struct{}

func NewHubHandler() *HubHandler {
	return &HubHandler{}
}

func (h *HubHandler) Page(w http.ResponseWriter, r *http.Request) {
	if err := hub.HubPage().Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

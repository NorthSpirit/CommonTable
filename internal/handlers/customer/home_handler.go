package customer

import (
	"net/http"

	"CommonTable/internal/templates"
)

type HomeHandler struct{}

func NewHomeHandler() *HomeHandler {
	return &HomeHandler{}
}

func (h *HomeHandler) Home(w http.ResponseWriter, r *http.Request) {
	err := templates.HomePage().Render(r.Context(), w)
	if err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

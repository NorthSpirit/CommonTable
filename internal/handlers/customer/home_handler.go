package customer

import (
	"net/http"

	"CommonTable/internal/services/homepageservices"
	"CommonTable/internal/templates"
)

type HomeHandler struct {
	homepageService *homepageservices.HomepageService
}

func NewHomeHandler(homepageService *homepageservices.HomepageService) *HomeHandler {
	return &HomeHandler{homepageService: homepageService}
}

func (h *HomeHandler) Home(w http.ResponseWriter, r *http.Request) {
	locale := getLocale(r)

	content, err := h.homepageService.GetContentForDisplay(r.Context(), locale)
	if err != nil {
		http.Error(w, "failed to load homepage content", http.StatusInternalServerError)
		return
	}

	view := templates.HomepageView{
		WelcomeText: content.WelcomeText,
	}
	if content.ImageOneUrl.Valid {
		view.ImageOneURL = content.ImageOneUrl.String
	}
	if content.ImageTwoUrl.Valid {
		view.ImageTwoURL = content.ImageTwoUrl.String
	}

	if err := templates.HomePage(view).Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

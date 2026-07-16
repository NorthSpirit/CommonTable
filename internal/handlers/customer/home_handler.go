package customer

import (
	"net/http"

	"CommonTable/internal/services/homepageservices"
	"CommonTable/internal/templates/client"
	"CommonTable/internal/templates/shared"
)

type HomeHandler struct {
	homepageService *homepageservices.HomepageService
}

func NewHomeHandler(homepageService *homepageservices.HomepageService) *HomeHandler {
	return &HomeHandler{homepageService: homepageService}
}

func (h *HomeHandler) Home(w http.ResponseWriter, r *http.Request) {
	locale := GetLocale(r)

	content, err := h.homepageService.GetContentForDisplay(r.Context(), locale)
	if err != nil {
		renderServerError(w, r)
		return
	}

	view := shared.HomepageView{
		WelcomeText: content.WelcomeText,
	}
	if content.ImageOneUrl.Valid {
		view.ImageOneURL = content.ImageOneUrl.String
	}
	if content.ImageTwoUrl.Valid {
		view.ImageTwoURL = content.ImageTwoUrl.String
	}

	if err := client.HomePage(locale, view).Render(r.Context(), w); err != nil {
		renderServerError(w, r)
	}
}

package customer

import (
	"net/http"

	"CommonTable/internal/services/homepageservices"
	"CommonTable/internal/services/restaurantservices"
	"CommonTable/internal/templates/client"
	"CommonTable/internal/templates/shared"
	"CommonTable/internal/textmarkup"
)

type HomeHandler struct {
	homepageService   *homepageservices.HomepageService
	restaurantService *restaurantservices.RestaurantService
}

func NewHomeHandler(homepageService *homepageservices.HomepageService, restaurantService *restaurantservices.RestaurantService) *HomeHandler {
	return &HomeHandler{homepageService: homepageService, restaurantService: restaurantService}
}

func (h *HomeHandler) Home(w http.ResponseWriter, r *http.Request) {
	locale := GetLocale(r)

	content, err := h.homepageService.GetContentForDisplay(r.Context(), locale)
	if err != nil {
		renderServerError(w, r)
		return
	}

	imageOneURL := ""
	if content.ImageOneUrl.Valid {
		imageOneURL = content.ImageOneUrl.String
	}
	imageTwoURL := ""
	if content.ImageTwoUrl.Valid {
		imageTwoURL = content.ImageTwoUrl.String
	}

	blocks := textmarkup.Parse(content.WelcomeText)
	blocks = textmarkup.EnsureImagesPlaced(blocks, imageOneURL != "", imageTwoURL != "")

	view := shared.HomepageView{
		WelcomeBlocks: blocks,
		ImageOneURL:   imageOneURL,
		ImageTwoURL:   imageTwoURL,
	}

	restaurants, err := h.restaurantService.GetRestaurantList(r.Context())
	if err != nil {
		renderServerError(w, r)
		return
	}

	restaurantViews := make([]shared.RestaurantView, 0, len(restaurants))
	for _, restaurant := range restaurants {
		restaurantViews = append(restaurantViews, toRestaurantView(restaurant))
	}

	if err := client.HomePage(locale, view, restaurantViews).Render(r.Context(), w); err != nil {
		renderServerError(w, r)
	}
}

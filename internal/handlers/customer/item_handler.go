package customer

import (
	"errors"
	"net/http"
	"strconv"

	"CommonTable/internal/helpers"
	"CommonTable/internal/services/foodrelationshipservices"
	"CommonTable/internal/services/foodservices"
	"CommonTable/internal/services/restaurantservices"
	"CommonTable/internal/templates/client"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

type ItemHandler struct {
	foodItemService       *foodservices.FoodItemService
	itemAllergenService   *foodrelationshipservices.ItemAllergenService
	ingredientItemService *foodrelationshipservices.IngredientItemService
	itemMenuService       *foodrelationshipservices.ItemMenuService
	restaurantService     *restaurantservices.RestaurantService
}

func NewItemHandler(
	foodItemService *foodservices.FoodItemService,
	itemAllergenService *foodrelationshipservices.ItemAllergenService,
	ingredientItemService *foodrelationshipservices.IngredientItemService,
	itemMenuService *foodrelationshipservices.ItemMenuService,
	restaurantService *restaurantservices.RestaurantService,
) *ItemHandler {
	return &ItemHandler{
		foodItemService:       foodItemService,
		itemAllergenService:   itemAllergenService,
		ingredientItemService: ingredientItemService,
		itemMenuService:       itemMenuService,
		restaurantService:     restaurantService,
	}
}

// Detail handles GET /items/{id}
func (h *ItemHandler) Detail(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		renderNotFound(w, r)
		return
	}

	locale := GetLocale(r)

	item, err := h.foodItemService.GetFoodItemForDisplay(r.Context(), int32(id), locale)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			renderNotFound(w, r)
			return
		}
		renderServerError(w, r)
		return
	}

	restaurant, err := h.restaurantService.GetRestaurantByID(r.Context(), item.RestaurantID)
	if err != nil {
		renderServerError(w, r)
		return
	}

	// GetFoodItemForDisplay omits price — fetch it separately from the
	// full item record rather than duplicating price into the display
	// query.
	fullItem, err := h.foodItemService.GetFoodItemByID(r.Context(), int32(id))
	if err != nil {
		renderServerError(w, r)
		return
	}
	price := 0.0
	if p, perr := helpers.NumericToFloat64(fullItem.Price); perr == nil && p != nil {
		price = *p
	}

	allergens, err := h.itemAllergenService.ListAllergensForItemDisplay(r.Context(), int32(id), locale)
	if err != nil {
		renderServerError(w, r)
		return
	}

	ingredients, err := h.ingredientItemService.ListIngredientsForItem(r.Context(), int32(id))
	if err != nil {
		renderServerError(w, r)
		return
	}

	// Which menu this item belongs to, for the back-link. An item can
	// technically sit on more than one menu — take the first as the
	// "came from" context, which is the common case.
	menus, err := h.itemMenuService.ListMenusForItem(r.Context(), int32(id))
	if err != nil {
		renderServerError(w, r)
		return
	}
	menuID := int32(0)
	if len(menus) > 0 {
		menuID = menus[0].ID
	}

	view := toItemDetailView(item, allergens, ingredients, price)
	branding := toBrandingView(restaurant)

	if err := client.ItemDetailPage(locale, branding, menuID, view).Render(r.Context(), w); err != nil {
		renderServerError(w, r)
	}
}

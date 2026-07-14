package customer

import (
	"errors"
	"net/http"
	"strconv"

	"CommonTable/internal/services/foodservices"
	"CommonTable/internal/templates"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

type MenuHandler struct {
	foodMenuService *foodservices.FoodMenuService
	foodItemService *foodservices.FoodItemService
}

func NewMenuHandler(
	foodMenuService *foodservices.FoodMenuService,
	foodItemService *foodservices.FoodItemService,
) *MenuHandler {
	return &MenuHandler{
		foodMenuService: foodMenuService,
		foodItemService: foodItemService,
	}
}

// Items handles GET /menus/{id}
func (h *MenuHandler) Items(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid menu id", http.StatusBadRequest)
		return
	}

	menu, err := h.foodMenuService.GetFoodMenuByID(r.Context(), int32(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "failed to load menu", http.StatusInternalServerError)
		return
	}

	items, err := h.foodItemService.ListMenuItemsForDisplay(r.Context(), int32(id))
	if err != nil {
		http.Error(w, "failed to load menu items", http.StatusInternalServerError)
		return
	}

	views := make([]templates.MenuItemView, 0, len(items))
	for _, item := range items {
		views = append(views, toMenuItemView(item))
	}

	if err := templates.MenuItemsPage(menu.Name, views).Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

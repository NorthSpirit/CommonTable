package chef

import (
	"errors"
	"net/http"
	"strconv"

	"CommonTable/internal/services/foodrelationshipservices"
	"CommonTable/internal/services/foodservices"
	"CommonTable/internal/services/restaurantservices"
	"CommonTable/internal/templates/chef"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

type MenuHandler struct {
	foodMenuService       *foodservices.FoodMenuService
	foodItemService       *foodservices.FoodItemService
	itemMenuService       *foodrelationshipservices.ItemMenuService
	itemAllergenService   *foodrelationshipservices.ItemAllergenService
	ingredientItemService *foodrelationshipservices.IngredientItemService
	restaurantService     *restaurantservices.RestaurantService
}

func NewMenuHandler(
	foodMenuService *foodservices.FoodMenuService,
	foodItemService *foodservices.FoodItemService,
	itemMenuService *foodrelationshipservices.ItemMenuService,
	itemAllergenService *foodrelationshipservices.ItemAllergenService,
	ingredientItemService *foodrelationshipservices.IngredientItemService,
	restaurantService *restaurantservices.RestaurantService,
) *MenuHandler {
	return &MenuHandler{
		foodMenuService:       foodMenuService,
		foodItemService:       foodItemService,
		itemMenuService:       itemMenuService,
		itemAllergenService:   itemAllergenService,
		ingredientItemService: ingredientItemService,
		restaurantService:     restaurantService,
	}
}

// List handles GET /chef/menus
func (h *MenuHandler) List(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromContext(r.Context())

	restaurant, err := h.restaurantService.GetRestaurantByID(r.Context(), claims.RestaurantID)
	if err != nil {
		http.Error(w, "failed to load restaurant", http.StatusInternalServerError)
		return
	}

	menus, err := h.foodMenuService.ListFoodMenusByRestaurant(r.Context(), claims.RestaurantID)
	if err != nil {
		http.Error(w, "failed to load menus", http.StatusInternalServerError)
		return
	}

	views := make([]chef.ChefMenuView, 0, len(menus))
	for _, m := range menus {
		view := chef.ChefMenuView{ID: m.ID, Name: m.Name}
		if m.Description.Valid {
			view.Description = m.Description.String
		}
		views = append(views, view)
	}

	if err := chef.MenusPage(restaurant.Name, claims.Name, views).Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

// MenuDetail handles GET /chef/menus/{menuID}
func (h *MenuHandler) MenuDetail(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromContext(r.Context())

	menuID, err := strconv.Atoi(chi.URLParam(r, "menuID"))
	if err != nil {
		http.Error(w, "invalid menu id", http.StatusBadRequest)
		return
	}

	menu, err := h.foodMenuService.GetFoodMenuByID(r.Context(), int32(menuID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "failed to load menu", http.StatusInternalServerError)
		return
	}

	items, err := h.itemMenuService.ListItemsInMenuForAdmin(r.Context(), int32(menuID))
	if err != nil {
		http.Error(w, "failed to load items", http.StatusInternalServerError)
		return
	}

	views := make([]chef.ChefItemListView, 0, len(items))
	for _, item := range items {
		if item.DeletedAt.Valid {
			continue
		}
		view := chef.ChefItemListView{ID: item.ID, Name: item.Name, IsAvailable: item.IsAvailable}
		if item.ImageUrl.Valid {
			view.ImageURL = item.ImageUrl.String
		}
		views = append(views, view)
	}

	if err := chef.ChefMenuDetailPage(int32(menuID), menu.Name, claims.Name, views).Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

// buildItemDetailView loads and assembles the full item detail view —
// shared by ItemDetail and UpdateInstructions so both stay in sync
// rather than maintaining two independent copies of this logic.
func (h *MenuHandler) buildItemDetailView(r *http.Request, menuID int32, itemID int32) (chef.ChefItemDetailView, error) {
	item, err := h.foodItemService.GetFoodItemForDisplay(r.Context(), itemID, "en")
	if err != nil {
		return chef.ChefItemDetailView{}, err
	}

	ingredients, err := h.ingredientItemService.ListIngredientsForItem(r.Context(), itemID)
	if err != nil {
		return chef.ChefItemDetailView{}, err
	}

	allergens, err := h.itemAllergenService.ListAllergensForItem(r.Context(), itemID)
	if err != nil {
		return chef.ChefItemDetailView{}, err
	}

	view := chef.ChefItemDetailView{
		ID:           item.ID,
		MenuID:       menuID,
		Name:         item.Name,
		Description:  item.Description,
		Instructions: item.Instructions,
		IsAvailable:  item.IsAvailable,
	}
	if item.ImageUrl.Valid {
		view.ImageURL = item.ImageUrl.String
	}
	for _, ing := range ingredients {
		ingView := chef.ChefIngredientView{Name: ing.Name}
		if ing.Measurement.Valid {
			ingView.Measurement = ing.Measurement.String
		}
		if ing.Amount.Valid {
			ingView.Amount = ing.Amount.Float64
		}
		view.Ingredients = append(view.Ingredients, ingView)
	}
	for _, a := range allergens {
		view.Allergens = append(view.Allergens, a.Name)
	}

	return view, nil
}

// ItemDetail handles GET /chef/menus/{menuID}/items/{itemID}
func (h *MenuHandler) ItemDetail(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromContext(r.Context())

	menuID, err := strconv.Atoi(chi.URLParam(r, "menuID"))
	if err != nil {
		http.Error(w, "invalid menu id", http.StatusBadRequest)
		return
	}
	itemID, err := strconv.Atoi(chi.URLParam(r, "itemID"))
	if err != nil {
		http.Error(w, "invalid item id", http.StatusBadRequest)
		return
	}

	view, err := h.buildItemDetailView(r, int32(menuID), int32(itemID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "failed to load item", http.StatusInternalServerError)
		return
	}

	if err := chef.ChefItemDetailPage(claims.Name, view, "").Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

// UpdateInstructions handles POST /chef/menus/{menuID}/items/{itemID}/instructions
func (h *MenuHandler) UpdateInstructions(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromContext(r.Context())

	menuID, err := strconv.Atoi(chi.URLParam(r, "menuID"))
	if err != nil {
		http.Error(w, "invalid menu id", http.StatusBadRequest)
		return
	}
	itemID, err := strconv.Atoi(chi.URLParam(r, "itemID"))
	if err != nil {
		http.Error(w, "invalid item id", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form submission", http.StatusBadRequest)
		return
	}

	instructions := r.FormValue("instructions")
	var instructionsPtr *string
	if instructions != "" {
		instructionsPtr = &instructions
	}

	if err := h.foodItemService.SetFoodItemInstructions(r.Context(), int32(itemID), instructionsPtr); err != nil {
		http.Error(w, "failed to update instructions", http.StatusInternalServerError)
		return
	}

	view, err := h.buildItemDetailView(r, int32(menuID), int32(itemID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "failed to load item", http.StatusInternalServerError)
		return
	}

	if err := chef.ChefItemDetailPage(claims.Name, view, "Instructions saved.").Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

// ToggleAvailability handles POST /chef/menus/{menuID}/items/{itemID}/availability
func (h *MenuHandler) ToggleAvailability(w http.ResponseWriter, r *http.Request) {
	menuID, err := strconv.Atoi(chi.URLParam(r, "menuID"))
	if err != nil {
		http.Error(w, "invalid menu id", http.StatusBadRequest)
		return
	}
	itemID, err := strconv.Atoi(chi.URLParam(r, "itemID"))
	if err != nil {
		http.Error(w, "invalid item id", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form submission", http.StatusBadRequest)
		return
	}

	available := r.FormValue("available") == "true"

	if err := h.foodItemService.SetFoodItemAvailability(r.Context(), int32(itemID), available); err != nil {
		http.Error(w, "failed to update availability", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/chef/menus/"+strconv.Itoa(menuID)+"/items/"+strconv.Itoa(itemID), http.StatusFound)
}

package manager

import (
	"errors"
	"net/http"
	"strconv"

	"CommonTable/internal/services/foodrelationshipservices"
	"CommonTable/internal/services/foodservices"
	"CommonTable/internal/services/restaurantservices"
	managertmpl "CommonTable/internal/templates/manager"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

type MenuHandler struct {
	foodMenuService   *foodservices.FoodMenuService
	itemMenuService   *foodrelationshipservices.ItemMenuService
	restaurantService *restaurantservices.RestaurantService
	uploadsDir        string
}

func NewMenuHandler(foodMenuService *foodservices.FoodMenuService, itemMenuService *foodrelationshipservices.ItemMenuService, restaurantService *restaurantservices.RestaurantService, uploadsDir string) *MenuHandler {
	return &MenuHandler{foodMenuService: foodMenuService, itemMenuService: itemMenuService, restaurantService: restaurantService, uploadsDir: uploadsDir}
}

// List handles GET /manager/menus
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

	views := make([]managertmpl.MenuView, 0, len(menus))
	for _, m := range menus {
		views = append(views, toMenuView(m))
	}

	if err := managertmpl.MenusPage(restaurant.Name, claims.Name, views).Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

// NewForm handles GET /manager/menus/new
func (h *MenuHandler) NewForm(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromContext(r.Context())
	if err := managertmpl.MenuFormPage(claims.Name, false, managertmpl.MenuView{}, "").Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

// Create handles POST /manager/menus/new
func (h *MenuHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromContext(r.Context())

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "invalid form submission", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	if name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}
	nameAlt := strPtrOrNil(r.FormValue("name_alt"))
	description := strPtrOrNil(r.FormValue("description"))
	descriptionAlt := strPtrOrNil(r.FormValue("description_alt"))

	dbMenu, err := h.foodMenuService.InsertFoodMenu(r.Context(), claims.RestaurantID, name, nameAlt, description, descriptionAlt, nil)
	if err != nil {
		http.Error(w, "failed to create menu", http.StatusInternalServerError)
		return
	}

	if files := r.MultipartForm.File["image"]; len(files) > 0 {
		url, err := saveUploadedImage(h.uploadsDir, "menus", strconv.Itoa(int(dbMenu.ID)), files[0])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		_ = h.foodMenuService.UpdateFoodMenu(r.Context(), dbMenu.ID, name, nameAlt, description, descriptionAlt, &url)
	}

	http.Redirect(w, r, "/manager/menus", http.StatusFound)
}

// EditForm handles GET /manager/menus/{menuID}/edit
func (h *MenuHandler) EditForm(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromContext(r.Context())

	menuID, err := strconv.Atoi(chi.URLParam(r, "menuID"))
	if err != nil {
		http.Error(w, "invalid menu id", http.StatusBadRequest)
		return
	}

	dbMenu, err := h.foodMenuService.GetFoodMenuByID(r.Context(), int32(menuID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "failed to load menu", http.StatusInternalServerError)
		return
	}
	if dbMenu.RestaurantID != claims.RestaurantID {
		http.NotFound(w, r)
		return
	}

	if err := managertmpl.MenuFormPage(claims.Name, true, toMenuView(dbMenu), "").Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

// Update handles POST /manager/menus/{menuID}/edit
func (h *MenuHandler) Update(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromContext(r.Context())

	menuID, err := strconv.Atoi(chi.URLParam(r, "menuID"))
	if err != nil {
		http.Error(w, "invalid menu id", http.StatusBadRequest)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "invalid form submission", http.StatusBadRequest)
		return
	}

	current, err := h.foodMenuService.GetFoodMenuByID(r.Context(), int32(menuID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "failed to load menu", http.StatusInternalServerError)
		return
	}
	if current.RestaurantID != claims.RestaurantID {
		http.NotFound(w, r)
		return
	}

	name := r.FormValue("name")
	if name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}
	nameAlt := strPtrOrNil(r.FormValue("name_alt"))
	description := strPtrOrNil(r.FormValue("description"))
	descriptionAlt := strPtrOrNil(r.FormValue("description_alt"))

	imageURL := ""
	if current.ImageUrl.Valid {
		imageURL = current.ImageUrl.String
	}
	if files := r.MultipartForm.File["image"]; len(files) > 0 {
		url, err := saveUploadedImage(h.uploadsDir, "menus", strconv.Itoa(menuID), files[0])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		imageURL = url
	}
	var imageURLPtr *string
	if imageURL != "" {
		imageURLPtr = &imageURL
	}

	if err := h.foodMenuService.UpdateFoodMenu(r.Context(), int32(menuID), name, nameAlt, description, descriptionAlt, imageURLPtr); err != nil {
		http.Error(w, "failed to update menu", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/manager/menus", http.StatusFound)
}

// Delete handles POST /manager/menus/{menuID}/delete
func (h *MenuHandler) Delete(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromContext(r.Context())

	menuID, err := strconv.Atoi(chi.URLParam(r, "menuID"))
	if err != nil {
		http.Error(w, "invalid menu id", http.StatusBadRequest)
		return
	}

	dbMenu, err := h.foodMenuService.GetFoodMenuByID(r.Context(), int32(menuID))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	if dbMenu.RestaurantID != claims.RestaurantID {
		http.NotFound(w, r)
		return
	}

	if err := h.foodMenuService.DeleteFoodMenu(r.Context(), int32(menuID)); err != nil {
		http.Error(w, "failed to delete menu", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/manager/menus", http.StatusFound)
}

// Detail handles GET /manager/menus/{menuID}
func (h *MenuHandler) Detail(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromContext(r.Context())

	menuID, err := strconv.Atoi(chi.URLParam(r, "menuID"))
	if err != nil {
		http.Error(w, "invalid menu id", http.StatusBadRequest)
		return
	}

	dbMenu, err := h.foodMenuService.GetFoodMenuByID(r.Context(), int32(menuID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "failed to load menu", http.StatusInternalServerError)
		return
	}
	if dbMenu.RestaurantID != claims.RestaurantID {
		http.NotFound(w, r)
		return
	}

	items, err := h.itemMenuService.ListItemsInMenuForAdmin(r.Context(), int32(menuID))
	if err != nil {
		http.Error(w, "failed to load items", http.StatusInternalServerError)
		return
	}

	itemViews := make([]managertmpl.FoodItemView, 0, len(items))
	for _, item := range items {
		itemViews = append(itemViews, toFoodItemViewFromMenuAdmin(item))
	}

	if err := managertmpl.MenuDetailPage(claims.Name, toMenuView(dbMenu), itemViews).Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

package admin

import (
	"errors"
	"net/http"
	"strconv"

	"CommonTable/internal/services/foodrelationshipservices"
	"CommonTable/internal/services/foodservices"
	"CommonTable/internal/services/restaurantservices"
	"CommonTable/internal/templates/admin/adminshared"
	"CommonTable/internal/templates/admin/menu"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

type MenuHandler struct {
	foodMenuService   *foodservices.FoodMenuService
	restaurantService *restaurantservices.RestaurantService
	itemMenuService   *foodrelationshipservices.ItemMenuService
	uploadsDir        string
}

func NewMenuHandler(foodMenuService *foodservices.FoodMenuService, restaurantService *restaurantservices.RestaurantService, itemMenuService *foodrelationshipservices.ItemMenuService, uploadsDir string) *MenuHandler {
	return &MenuHandler{foodMenuService: foodMenuService, restaurantService: restaurantService, itemMenuService: itemMenuService, uploadsDir: uploadsDir}
}

// List handles GET /admin/restaurants/{id}/menus
func (h *MenuHandler) List(w http.ResponseWriter, r *http.Request) {
	restaurantID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid restaurant id", http.StatusBadRequest)
		return
	}

	restaurant, err := h.restaurantService.GetRestaurantByID(r.Context(), int32(restaurantID))
	if err != nil {
		http.Error(w, "failed to load restaurant", http.StatusInternalServerError)
		return
	}

	menus, err := h.foodMenuService.ListFoodMenusByRestaurant(r.Context(), int32(restaurantID))
	if err != nil {
		http.Error(w, "failed to load menus", http.StatusInternalServerError)
		return
	}

	views := make([]adminshared.MenuView, 0, len(menus))
	for _, m := range menus {
		views = append(views, toMenuView(m))
	}

	if err := menu.MenusPage(int32(restaurantID), restaurant.Name, views).Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

// NewForm handles GET /admin/restaurants/{id}/menus/new
func (h *MenuHandler) NewForm(w http.ResponseWriter, r *http.Request) {
	restaurantID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid restaurant id", http.StatusBadRequest)
		return
	}
	if err := menu.MenuFormPage(int32(restaurantID), false, adminshared.MenuView{}, "").Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

// Create handles POST /admin/restaurants/{id}/menus/new
func (h *MenuHandler) Create(w http.ResponseWriter, r *http.Request) {
	restaurantID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid restaurant id", http.StatusBadRequest)
		return
	}
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

	dbMenu, err := h.foodMenuService.InsertFoodMenu(r.Context(), int32(restaurantID), name, nameAlt, description, descriptionAlt, nil)
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

	http.Redirect(w, r, "/admin/restaurants/"+strconv.Itoa(restaurantID)+"/menus", http.StatusFound)
}

// EditForm handles GET /admin/restaurants/{id}/menus/{menuID}/edit
func (h *MenuHandler) EditForm(w http.ResponseWriter, r *http.Request) {
	restaurantID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid restaurant id", http.StatusBadRequest)
		return
	}
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

	if err := menu.MenuFormPage(int32(restaurantID), true, toMenuView(dbMenu), "").Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

// Update handles POST /admin/restaurants/{id}/menus/{menuID}/edit
func (h *MenuHandler) Update(w http.ResponseWriter, r *http.Request) {
	restaurantID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid restaurant id", http.StatusBadRequest)
		return
	}
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

	http.Redirect(w, r, "/admin/restaurants/"+strconv.Itoa(restaurantID)+"/menus", http.StatusFound)
}

// Delete handles POST /admin/restaurants/{id}/menus/{menuID}/delete
func (h *MenuHandler) Delete(w http.ResponseWriter, r *http.Request) {
	restaurantID := chi.URLParam(r, "id")
	menuID, err := strconv.Atoi(chi.URLParam(r, "menuID"))
	if err != nil {
		http.Error(w, "invalid menu id", http.StatusBadRequest)
		return
	}
	if err := h.foodMenuService.DeleteFoodMenu(r.Context(), int32(menuID)); err != nil {
		http.Error(w, "failed to delete menu", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/admin/restaurants/"+restaurantID+"/menus", http.StatusFound)
}

// Detail handles GET /admin/restaurants/{id}/menus/{menuID}
func (h *MenuHandler) Detail(w http.ResponseWriter, r *http.Request) {
	restaurantID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid restaurant id", http.StatusBadRequest)
		return
	}
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

	items, err := h.itemMenuService.ListItemsInMenuForAdmin(r.Context(), int32(menuID))
	if err != nil {
		http.Error(w, "failed to load items", http.StatusInternalServerError)
		return
	}

	itemViews := make([]adminshared.FoodItemView, 0, len(items))
	for _, item := range items {
		itemViews = append(itemViews, toFoodItemViewFromMenuAdmin(item))
	}

	if err := menu.MenuDetailPage(int32(restaurantID), toMenuView(dbMenu), itemViews).Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

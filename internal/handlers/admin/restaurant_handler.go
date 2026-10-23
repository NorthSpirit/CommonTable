package admin

import (
	"errors"
	"net/http"
	"strconv"

	"CommonTable/internal/helpers"
	"CommonTable/internal/services/adminservices"
	"CommonTable/internal/services/restaurantservices"
	"CommonTable/internal/templates/admin/adminshared"
	"CommonTable/internal/templates/admin/restaurant"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

type RestaurantHandler struct {
	restaurantService *restaurantservices.RestaurantService
	adminService      *adminservices.AdminService
	uploadsDir        string
}

func NewRestaurantHandler(restaurantService *restaurantservices.RestaurantService, adminService *adminservices.AdminService, uploadsDir string) *RestaurantHandler {
	return &RestaurantHandler{restaurantService: restaurantService, adminService: adminService, uploadsDir: uploadsDir}
}

// List handles GET /admin/restaurants
func (h *RestaurantHandler) List(w http.ResponseWriter, r *http.Request) {
	restaurants, err := h.restaurantService.GetRestaurantListForAdmin(r.Context())
	if err != nil {
		http.Error(w, "failed to load restaurants", http.StatusInternalServerError)
		return
	}

	views := make([]adminshared.RestaurantView, 0, len(restaurants))
	for _, dbRestaurant := range restaurants {
		views = append(views, toRestaurantView(dbRestaurant))
	}

	if err := restaurant.RestaurantsPage(views).Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

// NewForm handles GET /admin/restaurants/new
func (h *RestaurantHandler) NewForm(w http.ResponseWriter, r *http.Request) {
	if err := restaurant.RestaurantFormPage(false, adminshared.RestaurantView{}, "").Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

// Create handles POST /admin/restaurants/new
func (h *RestaurantHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "invalid form submission", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	password := r.FormValue("password")

	view := adminshared.RestaurantView{
		Name:           name,
		Address:        r.FormValue("address"),
		Phone:          r.FormValue("phone"),
		Email:          r.FormValue("email"),
		Description:    r.FormValue("description"),
		DescriptionAlt: r.FormValue("description_alt"),
	}

	if name == "" || password == "" {
		h.renderFormError(w, r, false, view, "Name and password are required.")
		return
	}

	msg := helpers.CheckLengths(helpers.FieldLimit{Label: "Name", Value: name, Max: 255})
	if msg == "" {
		msg = helpers.ValidateContact(r.FormValue("email"), r.FormValue("phone"))
	}
	if msg == "" {
		msg = helpers.ValidateSecret("Password", password)
	}
	if msg != "" {
		h.renderFormError(w, r, false, view, msg)
		return
	}

	passwordHash, err := helpers.HashPassword(password)
	if err != nil {
		http.Error(w, "failed to hash password", http.StatusInternalServerError)
		return
	}

	address := strPtrOrNil(r.FormValue("address"))
	phone := strPtrOrNil(r.FormValue("phone"))
	email := strPtrOrNil(r.FormValue("email"))
	description := strPtrOrNil(r.FormValue("description"))
	descriptionAlt := strPtrOrNil(r.FormValue("description_alt"))

	dbRestaurant, err := h.restaurantService.InsertRestaurant(r.Context(), name, address, phone, email, passwordHash, nil, description, descriptionAlt)
	if err != nil {
		if isUniqueViolation(err) {
			h.renderFormError(w, r, false, view, "A restaurant named \""+name+"\" already exists.")
			return
		}
		if msg := helpers.PgUserMessage(err); msg != "" {
			h.renderFormError(w, r, false, view, msg)
			return
		}
		http.Error(w, "failed to create restaurant", http.StatusInternalServerError)
		return
	}

	if files := r.MultipartForm.File["logo"]; len(files) > 0 {
		url, err := saveUploadedImage(h.uploadsDir, "restaurants", strconv.Itoa(int(dbRestaurant.ID)), files[0])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := h.restaurantService.UpdateRestaurant(r.Context(), dbRestaurant.ID, name, address, phone, email, &url, description, descriptionAlt); err != nil {
			http.Error(w, "failed to save logo", http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, "/admin/restaurants", http.StatusFound)
}

// EditForm handles GET /admin/restaurants/{id}/edit
func (h *RestaurantHandler) EditForm(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid restaurant id", http.StatusBadRequest)
		return
	}

	dbRestaurant, err := h.restaurantService.GetRestaurantByID(r.Context(), int32(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "failed to load restaurant", http.StatusInternalServerError)
		return
	}

	view := toRestaurantView(dbRestaurant)

	if err := restaurant.RestaurantFormPage(true, view, "").Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

// Update handles POST /admin/restaurants/{id}/edit
func (h *RestaurantHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid restaurant id", http.StatusBadRequest)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "invalid form submission", http.StatusBadRequest)
		return
	}

	current, err := h.restaurantService.GetRestaurantByID(r.Context(), int32(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "failed to load restaurant", http.StatusInternalServerError)
		return
	}

	name := r.FormValue("name")
	password := r.FormValue("password")
	view := adminshared.RestaurantView{
		ID:             int32(id),
		Name:           name,
		Address:        r.FormValue("address"),
		Phone:          r.FormValue("phone"),
		Email:          r.FormValue("email"),
		Description:    r.FormValue("description"),
		DescriptionAlt: r.FormValue("description_alt"),
		LogoURL:        current.LogoUrl.String,
	}
	if name == "" {
		h.renderFormError(w, r, true, view, "Name is required.")
		return
	}

	msg := helpers.CheckLengths(helpers.FieldLimit{Label: "Name", Value: name, Max: 255})
	if msg == "" {
		msg = helpers.ValidateContact(r.FormValue("email"), r.FormValue("phone"))
	}
	if msg == "" {
		msg = helpers.ValidateSecret("Password", password)
	}
	if msg != "" {
		h.renderFormError(w, r, true, view, msg)
		return
	}

	address := strPtrOrNil(r.FormValue("address"))
	phone := strPtrOrNil(r.FormValue("phone"))
	email := strPtrOrNil(r.FormValue("email"))
	description := strPtrOrNil(r.FormValue("description"))
	descriptionAlt := strPtrOrNil(r.FormValue("description_alt"))

	logoURL := ""
	if current.LogoUrl.Valid {
		logoURL = current.LogoUrl.String
	}

	if files := r.MultipartForm.File["logo"]; len(files) > 0 {
		url, err := saveUploadedImage(h.uploadsDir, "restaurants", strconv.Itoa(id), files[0])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		logoURL = url
	}

	var logoURLPtr *string
	if logoURL != "" {
		logoURLPtr = &logoURL
	}

	if err := h.restaurantService.UpdateRestaurant(r.Context(), int32(id), name, address, phone, email, logoURLPtr, description, descriptionAlt); err != nil {
		if isUniqueViolation(err) {
			h.renderFormError(w, r, true, view, "A restaurant named \""+name+"\" already exists.")
			return
		}
		if msg := helpers.PgUserMessage(err); msg != "" {
			h.renderFormError(w, r, true, view, msg)
			return
		}
		http.Error(w, "failed to update restaurant", http.StatusInternalServerError)
		return
	}

	if password != "" {
		passwordHash, err := helpers.HashPassword(password)
		if err != nil {
			http.Error(w, "failed to hash password", http.StatusInternalServerError)
			return
		}
		if err := h.restaurantService.UpdateRestaurantPasswordHash(r.Context(), int32(id), passwordHash); err != nil {
			http.Error(w, "failed to update password", http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, "/admin/restaurants", http.StatusFound)
}

// DeleteConfirm handles GET /admin/restaurants/{id}/delete
func (h *RestaurantHandler) DeleteConfirm(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid restaurant id", http.StatusBadRequest)
		return
	}

	dbRestaurant, err := h.restaurantService.GetRestaurantByID(r.Context(), int32(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "failed to load restaurant", http.StatusInternalServerError)
		return
	}

	if err := restaurant.RestaurantDeleteConfirmPage(toRestaurantView(dbRestaurant), "").Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

// Delete handles POST /admin/restaurants/{id}/delete
func (h *RestaurantHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid restaurant id", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form submission", http.StatusBadRequest)
		return
	}

	password := r.FormValue("password")
	if err := h.adminService.VerifyPassword(password); err != nil {
		dbRestaurant, loadErr := h.restaurantService.GetRestaurantByID(r.Context(), int32(id))
		if loadErr != nil {
			http.Error(w, "failed to load restaurant", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		_ = restaurant.RestaurantDeleteConfirmPage(toRestaurantView(dbRestaurant), "Incorrect password.").Render(r.Context(), w)
		return
	}

	if err := h.restaurantService.DeleteRestaurant(r.Context(), int32(id)); err != nil {
		http.Error(w, "failed to delete restaurant", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/restaurants", http.StatusFound)
}

// Restore handles POST /admin/restaurants/{id}/restore
func (h *RestaurantHandler) Restore(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid restaurant id", http.StatusBadRequest)
		return
	}

	if err := h.restaurantService.RestoreRestaurant(r.Context(), int32(id)); err != nil {
		http.Error(w, "failed to restore restaurant", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/restaurants", http.StatusFound)
}

func (h *RestaurantHandler) renderFormError(w http.ResponseWriter, r *http.Request, isEdit bool, view adminshared.RestaurantView, message string) {
	w.WriteHeader(http.StatusBadRequest)
	_ = restaurant.RestaurantFormPage(isEdit, view, message).Render(r.Context(), w)
}

func strPtrOrNil(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// Detail handles GET /admin/restaurants/{id}
func (h *RestaurantHandler) Detail(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid restaurant id", http.StatusBadRequest)
		return
	}

	dbRestaurant, err := h.restaurantService.GetRestaurantByID(r.Context(), int32(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "failed to load restaurant", http.StatusInternalServerError)
		return
	}

	if err := restaurant.RestaurantDetailPage(toRestaurantView(dbRestaurant)).Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

package manager

import (
	"net/http"
	"strconv"

	"CommonTable/internal/helpers"
	"CommonTable/internal/repository"
	"CommonTable/internal/services/restaurantservices"
	managertmpl "CommonTable/internal/templates/manager"
)

type ProfileHandler struct {
	restaurantService *restaurantservices.RestaurantService
	uploadsDir        string
}

func NewProfileHandler(restaurantService *restaurantservices.RestaurantService, uploadsDir string) *ProfileHandler {
	return &ProfileHandler{restaurantService: restaurantService, uploadsDir: uploadsDir}
}

func buildProfileView(r repository.Restaurant) managertmpl.RestaurantProfileView {
	view := managertmpl.RestaurantProfileView{
		ID:   r.ID,
		Name: r.Name,
	}
	if r.Address.Valid {
		view.Address = r.Address.String
	}
	if r.Phone.Valid {
		view.Phone = r.Phone.String
	}
	if r.Email.Valid {
		view.Email = r.Email.String
	}
	if r.LogoUrl.Valid {
		view.LogoURL = r.LogoUrl.String
	}
	if r.Description.Valid {
		view.Description = r.Description.String
	}
	if r.DescriptionAlt.Valid {
		view.DescriptionAlt = r.DescriptionAlt.String
	}
	return view
}

// Edit handles GET /manager/profile
func (h *ProfileHandler) Edit(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromContext(r.Context())

	restaurant, err := h.restaurantService.GetRestaurantByID(r.Context(), claims.RestaurantID)
	if err != nil {
		http.Error(w, "failed to load restaurant", http.StatusInternalServerError)
		return
	}

	view := buildProfileView(restaurant)

	if err := managertmpl.RestaurantProfilePage(claims.Name, view, "", "").Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

// Update handles POST /manager/profile
func (h *ProfileHandler) Update(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromContext(r.Context())

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "invalid form submission", http.StatusBadRequest)
		return
	}

	current, err := h.restaurantService.GetRestaurantByID(r.Context(), claims.RestaurantID)
	if err != nil {
		http.Error(w, "failed to load restaurant", http.StatusInternalServerError)
		return
	}

	name := r.FormValue("name")
	if name == "" {
		view := buildProfileView(current)
		w.WriteHeader(http.StatusBadRequest)
		_ = managertmpl.RestaurantProfilePage(claims.Name, view, "Name is required.", "").Render(r.Context(), w)
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
		url, err := saveUploadedImage(h.uploadsDir, "restaurants", strconv.Itoa(int(claims.RestaurantID)), files[0])
		if err != nil {
			view := buildProfileView(current)
			w.WriteHeader(http.StatusBadRequest)
			_ = managertmpl.RestaurantProfilePage(claims.Name, view, err.Error(), "").Render(r.Context(), w)
			return
		}
		logoURL = url
	}

	var logoURLPtr *string
	if logoURL != "" {
		logoURLPtr = &logoURL
	}

	if err := h.restaurantService.UpdateRestaurant(r.Context(), claims.RestaurantID, name, address, phone, email, logoURLPtr, description, descriptionAlt); err != nil {
		http.Error(w, "failed to update restaurant", http.StatusInternalServerError)
		return
	}

	if newPassword := r.FormValue("password"); newPassword != "" {
		passwordHash, err := helpers.HashPassword(newPassword)
		if err != nil {
			http.Error(w, "failed to hash password", http.StatusInternalServerError)
			return
		}
		if err := h.restaurantService.UpdateRestaurantPasswordHash(r.Context(), claims.RestaurantID, passwordHash); err != nil {
			http.Error(w, "failed to update password", http.StatusInternalServerError)
			return
		}
	}

	updated, err := h.restaurantService.GetRestaurantByID(r.Context(), claims.RestaurantID)
	if err != nil {
		http.Error(w, "failed to reload restaurant", http.StatusInternalServerError)
		return
	}

	view := buildProfileView(updated)
	if err := managertmpl.RestaurantProfilePage(claims.Name, view, "", "Profile updated.").Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

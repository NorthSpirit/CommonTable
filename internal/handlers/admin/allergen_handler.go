package admin

import (
	"CommonTable/internal/services/foodrelationshipservices"
	"errors"
	"net/http"
	"strconv"

	"CommonTable/internal/services/foodservices"
	"CommonTable/internal/templates/admin/adminshared"
	"CommonTable/internal/templates/admin/allergen"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

type AllergenHandler struct {
	foodAllergenService *foodservices.FoodAllergenService
	itemAllergenService *foodrelationshipservices.ItemAllergenService
}

func NewAllergenHandler(foodAllergenService *foodservices.FoodAllergenService, itemAllergenService *foodrelationshipservices.ItemAllergenService) *AllergenHandler {
	return &AllergenHandler{foodAllergenService: foodAllergenService, itemAllergenService: itemAllergenService}
}

// List handles GET /admin/allergens
func (h *AllergenHandler) List(w http.ResponseWriter, r *http.Request) {
	allergens, err := h.foodAllergenService.GetFoodAllergensList(r.Context())
	if err != nil {
		http.Error(w, "failed to load allergens", http.StatusInternalServerError)
		return
	}

	views := make([]adminshared.AllergenAdminView, 0, len(allergens))
	for _, a := range allergens {
		views = append(views, toAllergenAdminView(a))
	}

	if err := allergen.AllergensAdminPage(views).Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

// NewForm handles GET /admin/allergens/new
func (h *AllergenHandler) NewForm(w http.ResponseWriter, r *http.Request) {
	if err := allergen.AllergenFormPage(false, adminshared.AllergenAdminView{}, "").Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

// Create handles POST /admin/allergens/new
func (h *AllergenHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form submission", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	code := r.FormValue("code")
	if name == "" || code == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = allergen.AllergenFormPage(false, adminshared.AllergenAdminView{Name: name, Code: code}, "Name and code are required.").Render(r.Context(), w)
		return
	}

	_, err := h.foodAllergenService.InsertFoodAllergen(
		r.Context(),
		name,
		strPtrOrNil(r.FormValue("name_alt")),
		strPtrOrNil(r.FormValue("description")),
		strPtrOrNil(r.FormValue("description_alt")),
		code,
	)
	if err != nil {
		http.Error(w, "failed to create allergen", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/allergens", http.StatusFound)
}

// EditForm handles GET /admin/allergens/{id}/edit
func (h *AllergenHandler) EditForm(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid allergen id", http.StatusBadRequest)
		return
	}

	a, err := h.foodAllergenService.GetFoodAllergenByID(r.Context(), int32(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "failed to load allergen", http.StatusInternalServerError)
		return
	}

	if err := allergen.AllergenFormPage(true, toAllergenAdminView(a), "").Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

// Update handles POST /admin/allergens/{id}/edit
func (h *AllergenHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid allergen id", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form submission", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	code := r.FormValue("code")
	if name == "" || code == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = allergen.AllergenFormPage(true, adminshared.AllergenAdminView{ID: int32(id), Name: name, Code: code}, "Name and code are required.").Render(r.Context(), w)
		return
	}

	err = h.foodAllergenService.UpdateFoodAllergen(
		r.Context(),
		int32(id),
		name,
		strPtrOrNil(r.FormValue("name_alt")),
		strPtrOrNil(r.FormValue("description")),
		strPtrOrNil(r.FormValue("description_alt")),
		code,
	)
	if err != nil {
		http.Error(w, "failed to update allergen", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/allergens", http.StatusFound)
}

// Delete handles POST /admin/allergens/{id}/delete
func (h *AllergenHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid allergen id", http.StatusBadRequest)
		return
	}

	if err := h.foodAllergenService.DeleteFoodAllergen(r.Context(), int32(id)); err != nil {
		http.Error(w, "failed to delete allergen", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/allergens", http.StatusFound)
}

// Usage handles GET /admin/allergens/{id}/usage
func (h *AllergenHandler) Usage(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid allergen id", http.StatusBadRequest)
		return
	}

	dbAllergen, err := h.foodAllergenService.GetFoodAllergenByID(r.Context(), int32(id))
	if err != nil {
		http.Error(w, "failed to load allergen", http.StatusInternalServerError)
		return
	}

	items, err := h.itemAllergenService.ListItemsWithAllergen(r.Context(), int32(id))
	if err != nil {
		http.Error(w, "failed to load usage", http.StatusInternalServerError)
		return
	}

	views := make([]adminshared.UsageItemView, 0, len(items))
	for _, item := range items {
		views = append(views, adminshared.UsageItemView{ID: item.ID, Name: item.Name})
	}

	if err := allergen.AllergenUsagePage(dbAllergen.Name, views).Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

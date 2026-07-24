package admin

import (
	"errors"
	"net/http"
	"strconv"

	"CommonTable/internal/services/foodrelationshipservices"
	"CommonTable/internal/services/foodservices"
	"CommonTable/internal/templates/admin/adminshared"
	"CommonTable/internal/templates/admin/ingredient"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

type IngredientHandler struct {
	foodIngredientService *foodservices.FoodIngredientService
	ingredientItemService *foodrelationshipservices.IngredientItemService
}

func NewIngredientHandler(foodIngredientService *foodservices.FoodIngredientService, ingredientItemService *foodrelationshipservices.IngredientItemService) *IngredientHandler {
	return &IngredientHandler{foodIngredientService: foodIngredientService, ingredientItemService: ingredientItemService}
}

// List handles GET /admin/ingredients
func (h *IngredientHandler) List(w http.ResponseWriter, r *http.Request) {
	ingredients, err := h.foodIngredientService.GetFoodIngredientList(r.Context())
	if err != nil {
		http.Error(w, "failed to load ingredients", http.StatusInternalServerError)
		return
	}

	views := make([]adminshared.IngredientView, 0, len(ingredients))
	for _, i := range ingredients {
		views = append(views, toIngredientView(i))
	}

	if err := ingredient.IngredientsPage(views).Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

// NewForm handles GET /admin/ingredients/new
func (h *IngredientHandler) NewForm(w http.ResponseWriter, r *http.Request) {
	if err := ingredient.IngredientFormPage(false, adminshared.IngredientView{}, "").Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

// Create handles POST /admin/ingredients/new
func (h *IngredientHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form submission", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = ingredient.IngredientFormPage(false, adminshared.IngredientView{}, "Name is required.").Render(r.Context(), w)
		return
	}

	_, err := h.foodIngredientService.InsertFoodIngredient(r.Context(), name, strPtrOrNil(r.FormValue("description")))
	if err != nil {
		http.Error(w, "failed to create ingredient", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/ingredients", http.StatusFound)
}

// EditForm handles GET /admin/ingredients/{id}/edit
func (h *IngredientHandler) EditForm(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid ingredient id", http.StatusBadRequest)
		return
	}

	foundIngredient, err := h.foodIngredientService.GetFoodIngredientByID(r.Context(), int32(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "failed to load ingredient", http.StatusInternalServerError)
		return
	}

	if err := ingredient.IngredientFormPage(true, toIngredientView(foundIngredient), "").Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

// Update handles POST /admin/ingredients/{id}/edit
func (h *IngredientHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid ingredient id", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form submission", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = ingredient.IngredientFormPage(true, adminshared.IngredientView{ID: int32(id)}, "Name is required.").Render(r.Context(), w)
		return
	}

	if err := h.foodIngredientService.UpdateFoodIngredient(r.Context(), int32(id), name, strPtrOrNil(r.FormValue("description"))); err != nil {
		http.Error(w, "failed to update ingredient", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/ingredients", http.StatusFound)
}

// Delete handles POST /admin/ingredients/{id}/delete
func (h *IngredientHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid ingredient id", http.StatusBadRequest)
		return
	}

	if err := h.foodIngredientService.DeleteFoodIngredient(r.Context(), int32(id)); err != nil {
		http.Error(w, "failed to delete ingredient", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/ingredients", http.StatusFound)
}

// Usage handles GET /admin/ingredients/{id}/usage
func (h *IngredientHandler) Usage(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid ingredient id", http.StatusBadRequest)
		return
	}

	dbIngredient, err := h.foodIngredientService.GetFoodIngredientByID(r.Context(), int32(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "failed to load ingredient", http.StatusInternalServerError)
		return
	}

	items, err := h.ingredientItemService.ListItemsUsingIngredient(r.Context(), int32(id))
	if err != nil {
		http.Error(w, "failed to load usage", http.StatusInternalServerError)
		return
	}

	views := make([]adminshared.UsageItemView, 0, len(items))
	for _, item := range items {
		views = append(views, adminshared.UsageItemView{ID: item.ID, Name: item.Name})
	}

	if err := ingredient.IngredientUsagePage(dbIngredient.Name, views).Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

package admin

import (
	"errors"
	"net/http"
	"strconv"

	"CommonTable/internal/helpers"
	"CommonTable/internal/services/foodrelationshipservices"
	"CommonTable/internal/services/foodservices"
	"CommonTable/internal/templates/admin/adminshared"
	"CommonTable/internal/templates/admin/fooditem"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

type FoodItemHandler struct {
	foodItemService       *foodservices.FoodItemService
	foodMenuService       *foodservices.FoodMenuService
	foodAllergenService   *foodservices.FoodAllergenService
	foodIngredientService *foodservices.FoodIngredientService
	itemAllergenService   *foodrelationshipservices.ItemAllergenService
	itemMenuService       *foodrelationshipservices.ItemMenuService
	ingredientItemService *foodrelationshipservices.IngredientItemService
	uploadsDir            string
}

func NewFoodItemHandler(
	foodItemService *foodservices.FoodItemService,
	foodMenuService *foodservices.FoodMenuService,
	foodAllergenService *foodservices.FoodAllergenService,
	foodIngredientService *foodservices.FoodIngredientService,
	itemAllergenService *foodrelationshipservices.ItemAllergenService,
	itemMenuService *foodrelationshipservices.ItemMenuService,
	ingredientItemService *foodrelationshipservices.IngredientItemService,
	uploadsDir string,
) *FoodItemHandler {
	return &FoodItemHandler{
		foodItemService:       foodItemService,
		foodMenuService:       foodMenuService,
		foodAllergenService:   foodAllergenService,
		foodIngredientService: foodIngredientService,
		itemAllergenService:   itemAllergenService,
		itemMenuService:       itemMenuService,
		ingredientItemService: ingredientItemService,
		uploadsDir:            uploadsDir,
	}
}

func (h *FoodItemHandler) buildAllergenCheckboxes(r *http.Request, itemID int32) ([]adminshared.AllergenCheckboxView, error) {
	allAllergens, err := h.foodAllergenService.GetFoodAllergensList(r.Context())
	if err != nil {
		return nil, err
	}

	checked := map[int32]bool{}
	if itemID != 0 {
		linked, err := h.itemAllergenService.ListAllergensForItem(r.Context(), itemID)
		if err != nil {
			return nil, err
		}
		for _, a := range linked {
			checked[a.ID] = true
		}
	}

	views := make([]adminshared.AllergenCheckboxView, 0, len(allAllergens))
	for _, a := range allAllergens {
		views = append(views, adminshared.AllergenCheckboxView{ID: a.ID, Name: a.Name, Checked: checked[a.ID]})
	}
	return views, nil
}

func (h *FoodItemHandler) buildIngredientOptions(r *http.Request) ([]adminshared.IngredientOptionView, error) {
	all, err := h.foodIngredientService.GetFoodIngredientList(r.Context())
	if err != nil {
		return nil, err
	}
	views := make([]adminshared.IngredientOptionView, 0, len(all))
	for _, ing := range all {
		views = append(views, adminshared.IngredientOptionView{ID: ing.ID, Name: ing.Name})
	}
	return views, nil
}

func (h *FoodItemHandler) buildItemIngredients(r *http.Request, itemID int32) ([]adminshared.ItemIngredientView, error) {
	linked, err := h.ingredientItemService.ListIngredientsForItem(r.Context(), itemID)
	if err != nil {
		return nil, err
	}
	views := make([]adminshared.ItemIngredientView, 0, len(linked))
	for _, ing := range linked {
		view := adminshared.ItemIngredientView{IngredientID: ing.ID, Name: ing.Name}
		if ing.Measurement.Valid {
			view.Measurement = ing.Measurement.String
		}
		if ing.Amount.Valid {
			view.Amount = ing.Amount.Float64
		}
		views = append(views, view)
	}
	return views, nil
}

func (h *FoodItemHandler) buildMenuCheckboxes(r *http.Request, restaurantID int32, itemID int32) ([]adminshared.MenuCheckboxView, error) {
	allMenus, err := h.foodMenuService.ListFoodMenusByRestaurant(r.Context(), restaurantID)
	if err != nil {
		return nil, err
	}

	checked := map[int32]bool{}
	if itemID != 0 {
		linked, err := h.itemMenuService.ListMenusForItem(r.Context(), itemID)
		if err != nil {
			return nil, err
		}
		for _, m := range linked {
			checked[m.ID] = true
		}
	}

	views := make([]adminshared.MenuCheckboxView, 0, len(allMenus))
	for _, m := range allMenus {
		views = append(views, adminshared.MenuCheckboxView{ID: m.ID, Name: m.Name, Checked: checked[m.ID]})
	}
	return views, nil
}

// buildMenuCheckboxesWithDefault pre-checks defaultMenuID even for a
// brand-new item, so the "new item" form starts with the menu the admin
// came from already checked.
func (h *FoodItemHandler) buildMenuCheckboxesWithDefault(r *http.Request, restaurantID int32, itemID int32, defaultMenuID int32) ([]adminshared.MenuCheckboxView, error) {
	views, err := h.buildMenuCheckboxes(r, restaurantID, itemID)
	if err != nil {
		return nil, err
	}
	for i := range views {
		if views[i].ID == defaultMenuID {
			views[i].Checked = true
		}
	}
	return views, nil
}

// NewForm handles GET /admin/restaurants/{id}/menus/{menuID}/items/new
func (h *FoodItemHandler) NewForm(w http.ResponseWriter, r *http.Request) {
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

	allergens, err := h.buildAllergenCheckboxes(r, 0)
	if err != nil {
		http.Error(w, "failed to load allergens", http.StatusInternalServerError)
		return
	}

	menus, err := h.buildMenuCheckboxesWithDefault(r, int32(restaurantID), 0, int32(menuID))
	if err != nil {
		http.Error(w, "failed to load menus", http.StatusInternalServerError)
		return
	}

	if err := fooditem.FoodItemFormPage(int32(restaurantID), int32(menuID), false, adminshared.FoodItemView{}, allergens, menus, nil, "").Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

// Create handles POST /admin/restaurants/{id}/menus/{menuID}/items/new
func (h *FoodItemHandler) Create(w http.ResponseWriter, r *http.Request) {
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

	name := r.FormValue("name")
	if name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	if msg := helpers.ValidateFoodItemForm(r); msg != "" {
		h.renderFormError(w, r, int32(restaurantID), int32(menuID), 0, false, "", msg)
		return
	}

	price, _ := strconv.ParseFloat(r.FormValue("price"), 64)
	prepTime := parseIntPtr(r.FormValue("prep_time_mins"))
	isAvailable := r.FormValue("is_available") == "on"

	nameAlt := strPtrOrNil(r.FormValue("name_alt"))
	description := strPtrOrNil(r.FormValue("description"))
	descriptionAlt := strPtrOrNil(r.FormValue("description_alt"))
	instructions := strPtrOrNil(r.FormValue("instructions"))
	instructionsAlt := strPtrOrNil(r.FormValue("instructions_alt"))

	item, err := h.foodItemService.InsertFoodItem(
		r.Context(),
		int32(restaurantID),
		name,
		nameAlt,
		description,
		descriptionAlt,
		isAvailable,
		nil,
		prepTime,
		instructions,
		instructionsAlt,
		price,
	)
	if err != nil {
		if msg := helpers.PgUserMessage(err); msg != "" {
			h.renderFormError(w, r, int32(restaurantID), int32(menuID), 0, false, "", msg)
			return
		}
		http.Error(w, "failed to create food item", http.StatusInternalServerError)
		return
	}

	if files := r.MultipartForm.File["image"]; len(files) > 0 {
		url, err := saveUploadedImage(h.uploadsDir, "items", strconv.Itoa(int(item.ID)), files[0])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		_ = h.foodItemService.UpdateFoodItem(r.Context(), item.ID, name, nameAlt, description, descriptionAlt, isAvailable, &url, prepTime, instructions, instructionsAlt, price)
	}

	// Always assign to the menu this item was created from.
	if err := h.itemMenuService.AddFoodItemToMenu(r.Context(), int32(menuID), item.ID); err != nil {
		http.Error(w, "failed to assign item to menu", http.StatusInternalServerError)
		return
	}
	// Plus anything else checked on the form.
	for _, raw := range r.Form["menus"] {
		id, err := strconv.Atoi(raw)
		if err != nil || int32(id) == int32(menuID) {
			continue
		}
		_ = h.itemMenuService.AddFoodItemToMenu(r.Context(), int32(id), item.ID)
	}

	http.Redirect(w, r, "/admin/restaurants/"+strconv.Itoa(restaurantID)+"/menus/"+strconv.Itoa(menuID), http.StatusFound)
}

// EditForm handles GET /admin/restaurants/{id}/menus/{menuID}/items/{itemID}/edit
func (h *FoodItemHandler) EditForm(w http.ResponseWriter, r *http.Request) {
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
	itemID, err := strconv.Atoi(chi.URLParam(r, "itemID"))
	if err != nil {
		http.Error(w, "invalid item id", http.StatusBadRequest)
		return
	}

	item, err := h.foodItemService.GetFoodItemByID(r.Context(), int32(itemID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "failed to load food item", http.StatusInternalServerError)
		return
	}

	allergens, err := h.buildAllergenCheckboxes(r, int32(itemID))
	if err != nil {
		http.Error(w, "failed to load allergens", http.StatusInternalServerError)
		return
	}

	menus, err := h.buildMenuCheckboxes(r, int32(restaurantID), int32(itemID))
	if err != nil {
		http.Error(w, "failed to load menus", http.StatusInternalServerError)
		return
	}

	ingredientOptions, err := h.buildIngredientOptions(r)
	if err != nil {
		http.Error(w, "failed to load ingredients", http.StatusInternalServerError)
		return
	}

	view := toFoodItemViewFromGet(item)
	view.Ingredients, err = h.buildItemIngredients(r, int32(itemID))
	if err != nil {
		http.Error(w, "failed to load item ingredients", http.StatusInternalServerError)
		return
	}

	if err := fooditem.FoodItemFormPage(int32(restaurantID), int32(menuID), true, view, allergens, menus, ingredientOptions, "").Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

// Update handles POST /admin/restaurants/{id}/menus/{menuID}/items/{itemID}/edit
func (h *FoodItemHandler) Update(w http.ResponseWriter, r *http.Request) {
	restaurantID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid restaurant id", http.StatusBadRequest)
		return
	}
	menuID := chi.URLParam(r, "menuID")
	itemID, err := strconv.Atoi(chi.URLParam(r, "itemID"))
	if err != nil {
		http.Error(w, "invalid item id", http.StatusBadRequest)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "invalid form submission", http.StatusBadRequest)
		return
	}

	// Ingredient add/remove are handled as their own quick round-trip,
	// separate from the main item Save — bypasses the required-field
	// validation on name/price via formnovalidate on those buttons.
	if r.FormValue("add_ingredient") != "" {
		ingredientID, errID := strconv.Atoi(r.FormValue("new_ingredient_id"))
		amount, _ := strconv.ParseFloat(r.FormValue("new_ingredient_amount"), 64)
		measurement := r.FormValue("new_ingredient_measurement")
		if errID == nil && ingredientID != 0 {
			_ = h.ingredientItemService.AddIngredientToItem(r.Context(), int32(ingredientID), int32(itemID), strPtrOrNil(measurement), &amount)
		}
		http.Redirect(w, r, "/admin/restaurants/"+strconv.Itoa(restaurantID)+"/menus/"+menuID+"/items/"+strconv.Itoa(itemID)+"/edit", http.StatusFound)
		return
	}

	if removeID := r.FormValue("remove_ingredient"); removeID != "" {
		ingredientID, errID := strconv.Atoi(removeID)
		if errID == nil {
			_ = h.ingredientItemService.RemoveIngredientFromItem(r.Context(), int32(ingredientID), int32(itemID))
		}
		http.Redirect(w, r, "/admin/restaurants/"+strconv.Itoa(restaurantID)+"/menus/"+menuID+"/items/"+strconv.Itoa(itemID)+"/edit", http.StatusFound)
		return
	}

	current, err := h.foodItemService.GetFoodItemByID(r.Context(), int32(itemID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "failed to load food item", http.StatusInternalServerError)
		return
	}

	name := r.FormValue("name")
	if name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	if msg := helpers.ValidateFoodItemForm(r); msg != "" {
		menuIDNum, _ := strconv.Atoi(menuID)
		h.renderFormError(w, r, int32(restaurantID), int32(menuIDNum), int32(itemID), true, current.ImageUrl.String, msg)
		return
	}

	price, _ := strconv.ParseFloat(r.FormValue("price"), 64)
	prepTime := parseIntPtr(r.FormValue("prep_time_mins"))
	isAvailable := r.FormValue("is_available") == "on"

	nameAlt := strPtrOrNil(r.FormValue("name_alt"))
	description := strPtrOrNil(r.FormValue("description"))
	descriptionAlt := strPtrOrNil(r.FormValue("description_alt"))
	instructions := strPtrOrNil(r.FormValue("instructions"))
	instructionsAlt := strPtrOrNil(r.FormValue("instructions_alt"))

	imageURL := ""
	if current.ImageUrl.Valid {
		imageURL = current.ImageUrl.String
	}

	if files := r.MultipartForm.File["image"]; len(files) > 0 {
		url, err := saveUploadedImage(h.uploadsDir, "items", strconv.Itoa(itemID), files[0])
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

	err = h.foodItemService.UpdateFoodItem(
		r.Context(),
		int32(itemID),
		name,
		nameAlt,
		description,
		descriptionAlt,
		isAvailable,
		imageURLPtr,
		prepTime,
		instructions,
		instructionsAlt,
		price,
	)
	if err != nil {
		if msg := helpers.PgUserMessage(err); msg != "" {
			menuIDNum, _ := strconv.Atoi(menuID)
			h.renderFormError(w, r, int32(restaurantID), int32(menuIDNum), int32(itemID), true, current.ImageUrl.String, msg)
			return
		}
		http.Error(w, "failed to update food item", http.StatusInternalServerError)
		return
	}

	if err := h.syncAllergens(r, int32(itemID)); err != nil {
		http.Error(w, "failed to update allergens", http.StatusInternalServerError)
		return
	}

	if err := h.syncMenus(r, int32(itemID)); err != nil {
		http.Error(w, "failed to update menus", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/restaurants/"+strconv.Itoa(restaurantID)+"/menus/"+menuID, http.StatusFound)
}

func (h *FoodItemHandler) syncAllergens(r *http.Request, itemID int32) error {
	checkedIDs := map[int32]bool{}
	for _, raw := range r.Form["allergens"] {
		if id, err := strconv.Atoi(raw); err == nil {
			checkedIDs[int32(id)] = true
		}
	}

	current, err := h.itemAllergenService.ListAllergensForItem(r.Context(), itemID)
	if err != nil {
		return err
	}
	currentIDs := map[int32]bool{}
	for _, a := range current {
		currentIDs[a.ID] = true
	}

	for id := range checkedIDs {
		if !currentIDs[id] {
			if err := h.itemAllergenService.AddAllergenToItem(r.Context(), itemID, id); err != nil {
				return err
			}
		}
	}
	for id := range currentIDs {
		if !checkedIDs[id] {
			if err := h.itemAllergenService.RemoveAllergenFromItem(r.Context(), itemID, id); err != nil {
				return err
			}
		}
	}
	return nil
}

func (h *FoodItemHandler) syncMenus(r *http.Request, itemID int32) error {
	checkedIDs := map[int32]bool{}
	for _, raw := range r.Form["menus"] {
		if id, err := strconv.Atoi(raw); err == nil {
			checkedIDs[int32(id)] = true
		}
	}

	current, err := h.itemMenuService.ListMenusForItem(r.Context(), itemID)
	if err != nil {
		return err
	}
	currentIDs := map[int32]bool{}
	for _, m := range current {
		currentIDs[m.ID] = true
	}

	for id := range checkedIDs {
		if !currentIDs[id] {
			if err := h.itemMenuService.AddFoodItemToMenu(r.Context(), id, itemID); err != nil {
				return err
			}
		}
	}
	for id := range currentIDs {
		if !checkedIDs[id] {
			if err := h.itemMenuService.RemoveFoodItemFromMenu(r.Context(), id, itemID); err != nil {
				return err
			}
		}
	}
	return nil
}

// Delete handles POST /admin/restaurants/{id}/menus/{menuID}/items/{itemID}/delete
func (h *FoodItemHandler) Delete(w http.ResponseWriter, r *http.Request) {
	restaurantID := chi.URLParam(r, "id")
	menuID := chi.URLParam(r, "menuID")
	itemID, err := strconv.Atoi(chi.URLParam(r, "itemID"))
	if err != nil {
		http.Error(w, "invalid item id", http.StatusBadRequest)
		return
	}
	if err := h.foodItemService.DeleteFoodItem(r.Context(), int32(itemID)); err != nil {
		http.Error(w, "failed to delete food item", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/admin/restaurants/"+restaurantID+"/menus/"+menuID, http.StatusFound)
}

// Restore handles POST /admin/restaurants/{id}/menus/{menuID}/items/{itemID}/restore
func (h *FoodItemHandler) Restore(w http.ResponseWriter, r *http.Request) {
	restaurantID := chi.URLParam(r, "id")
	menuID := chi.URLParam(r, "menuID")
	itemID, err := strconv.Atoi(chi.URLParam(r, "itemID"))
	if err != nil {
		http.Error(w, "invalid item id", http.StatusBadRequest)
		return
	}
	if err := h.foodItemService.RestoreFoodItem(r.Context(), int32(itemID)); err != nil {
		http.Error(w, "failed to restore food item", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/admin/restaurants/"+restaurantID+"/menus/"+menuID, http.StatusFound)
}

func parseIntPtr(s string) *int32 {
	if s == "" {
		return nil
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return nil
	}
	v32 := int32(v)
	return &v32
}

// renderFormError re-renders the item form with the submitted values and an
// error message. Allergen / menu checkboxes are reloaded from the database.
func (h *FoodItemHandler) renderFormError(w http.ResponseWriter, r *http.Request, restaurantID, menuID, itemID int32, isEdit bool, imageURL, msg string) {
	view := adminshared.FoodItemView{
		ID:              itemID,
		Name:            r.FormValue("name"),
		NameAlt:         r.FormValue("name_alt"),
		Description:     r.FormValue("description"),
		DescriptionAlt:  r.FormValue("description_alt"),
		Instructions:    r.FormValue("instructions"),
		InstructionsAlt: r.FormValue("instructions_alt"),
		Price:           helpers.SafeFormPrice(r),
		ImageURL:        imageURL,
		PrepTimeMins:    helpers.SafeFormPrepTime(r),
		IsAvailable:     r.FormValue("is_available") == "on",
	}

	allergens, err := h.buildAllergenCheckboxes(r, itemID)
	if err != nil {
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	var menus []adminshared.MenuCheckboxView
	var ingredientOptions []adminshared.IngredientOptionView
	if isEdit {
		menus, err = h.buildMenuCheckboxes(r, restaurantID, itemID)
		if err == nil {
			ingredientOptions, err = h.buildIngredientOptions(r)
		}
		if err == nil {
			view.Ingredients, err = h.buildItemIngredients(r, itemID)
		}
	} else {
		menus, err = h.buildMenuCheckboxesWithDefault(r, restaurantID, 0, menuID)
	}
	if err != nil {
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusBadRequest)
	_ = fooditem.FoodItemFormPage(restaurantID, menuID, isEdit, view, allergens, menus, ingredientOptions, msg).Render(r.Context(), w)
}

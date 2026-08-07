package manager

import (
	"errors"
	"net/http"
	"strconv"

	"CommonTable/internal/services/foodrelationshipservices"
	"CommonTable/internal/services/foodservices"
	managertmpl "CommonTable/internal/templates/manager"

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

func (h *FoodItemHandler) buildAllergenCheckboxes(r *http.Request, itemID int32) ([]managertmpl.AllergenCheckboxView, error) {
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

	views := make([]managertmpl.AllergenCheckboxView, 0, len(allAllergens))
	for _, a := range allAllergens {
		views = append(views, managertmpl.AllergenCheckboxView{ID: a.ID, Name: a.Name, Checked: checked[a.ID]})
	}
	return views, nil
}

func (h *FoodItemHandler) buildIngredientOptions(r *http.Request) ([]managertmpl.IngredientOptionView, error) {
	all, err := h.foodIngredientService.GetFoodIngredientList(r.Context())
	if err != nil {
		return nil, err
	}
	views := make([]managertmpl.IngredientOptionView, 0, len(all))
	for _, ing := range all {
		views = append(views, toIngredientOptionView(ing))
	}
	return views, nil
}

func (h *FoodItemHandler) buildItemIngredients(r *http.Request, itemID int32) ([]managertmpl.ItemIngredientView, error) {
	linked, err := h.ingredientItemService.ListIngredientsForItem(r.Context(), itemID)
	if err != nil {
		return nil, err
	}
	views := make([]managertmpl.ItemIngredientView, 0, len(linked))
	for _, ing := range linked {
		views = append(views, toItemIngredientView(ing))
	}
	return views, nil
}

// verifyMenuOwnership confirms the given menu belongs to the manager's
// own restaurant — without this, a manager could add items to another
// restaurant's menu just by guessing/incrementing the menuID in the URL.
func (h *FoodItemHandler) verifyMenuOwnership(r *http.Request, menuID int32, restaurantID int32) error {
	menu, err := h.foodMenuService.GetFoodMenuByID(r.Context(), menuID)
	if err != nil {
		return err
	}
	if menu.RestaurantID != restaurantID {
		return pgx.ErrNoRows
	}
	return nil
}

// NewForm handles GET /manager/menus/{menuID}/items/new
func (h *FoodItemHandler) NewForm(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromContext(r.Context())

	menuID, err := strconv.Atoi(chi.URLParam(r, "menuID"))
	if err != nil {
		http.Error(w, "invalid menu id", http.StatusBadRequest)
		return
	}
	if err := h.verifyMenuOwnership(r, int32(menuID), claims.RestaurantID); err != nil {
		http.NotFound(w, r)
		return
	}

	allergens, err := h.buildAllergenCheckboxes(r, 0)
	if err != nil {
		http.Error(w, "failed to load allergens", http.StatusInternalServerError)
		return
	}

	if err := managertmpl.FoodItemFormPage(claims.Name, int32(menuID), false, managertmpl.FoodItemView{}, allergens, nil, nil, "").Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

// Create handles POST /manager/menus/{menuID}/items/new
func (h *FoodItemHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromContext(r.Context())

	menuID, err := strconv.Atoi(chi.URLParam(r, "menuID"))
	if err != nil {
		http.Error(w, "invalid menu id", http.StatusBadRequest)
		return
	}
	if err := h.verifyMenuOwnership(r, int32(menuID), claims.RestaurantID); err != nil {
		http.NotFound(w, r)
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
		claims.RestaurantID,
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

	if err := h.itemMenuService.AddFoodItemToMenu(r.Context(), int32(menuID), item.ID); err != nil {
		http.Error(w, "failed to assign item to menu", http.StatusInternalServerError)
		return
	}

	if err := h.syncAllergens(r, item.ID); err != nil {
		http.Error(w, "failed to assign allergens", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/manager/menus/"+strconv.Itoa(menuID)+"/items/"+strconv.Itoa(int(item.ID))+"/edit", http.StatusFound)
}

// EditForm handles GET /manager/menus/{menuID}/items/{itemID}/edit
func (h *FoodItemHandler) EditForm(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromContext(r.Context())

	menuID, err := strconv.Atoi(chi.URLParam(r, "menuID"))
	if err != nil {
		http.Error(w, "invalid menu id", http.StatusBadRequest)
		return
	}
	if err := h.verifyMenuOwnership(r, int32(menuID), claims.RestaurantID); err != nil {
		http.NotFound(w, r)
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

	ingredientOptions, err := h.buildIngredientOptions(r)
	if err != nil {
		http.Error(w, "failed to load ingredients", http.StatusInternalServerError)
		return
	}

	itemIngredients, err := h.buildItemIngredients(r, int32(itemID))
	if err != nil {
		http.Error(w, "failed to load item ingredients", http.StatusInternalServerError)
		return
	}

	view := toFoodItemViewFromGet(item)

	if err := managertmpl.FoodItemFormPage(claims.Name, int32(menuID), true, view, allergens, ingredientOptions, itemIngredients, "").Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

// Update handles POST /manager/menus/{menuID}/items/{itemID}/edit
func (h *FoodItemHandler) Update(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromContext(r.Context())

	menuID, err := strconv.Atoi(chi.URLParam(r, "menuID"))
	if err != nil {
		http.Error(w, "invalid menu id", http.StatusBadRequest)
		return
	}
	if err := h.verifyMenuOwnership(r, int32(menuID), claims.RestaurantID); err != nil {
		http.NotFound(w, r)
		return
	}

	itemID, err := strconv.Atoi(chi.URLParam(r, "itemID"))
	if err != nil {
		http.Error(w, "invalid item id", http.StatusBadRequest)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "invalid form submission", http.StatusBadRequest)
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
		http.Error(w, "failed to update food item", http.StatusInternalServerError)
		return
	}

	if err := h.syncAllergens(r, int32(itemID)); err != nil {
		http.Error(w, "failed to update allergens", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/manager/menus/"+strconv.Itoa(menuID)+"/items/"+strconv.Itoa(itemID)+"/edit", http.StatusFound)
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

// Delete handles POST /manager/menus/{menuID}/items/{itemID}/delete
func (h *FoodItemHandler) Delete(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromContext(r.Context())

	menuID := chi.URLParam(r, "menuID")
	menuIDInt, err := strconv.Atoi(menuID)
	if err != nil {
		http.Error(w, "invalid menu id", http.StatusBadRequest)
		return
	}
	if err := h.verifyMenuOwnership(r, int32(menuIDInt), claims.RestaurantID); err != nil {
		http.NotFound(w, r)
		return
	}

	itemID, err := strconv.Atoi(chi.URLParam(r, "itemID"))
	if err != nil {
		http.Error(w, "invalid item id", http.StatusBadRequest)
		return
	}
	if err := h.foodItemService.DeleteFoodItem(r.Context(), int32(itemID)); err != nil {
		http.Error(w, "failed to delete food item", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/manager/menus/"+menuID, http.StatusFound)
}

// Restore handles POST /manager/menus/{menuID}/items/{itemID}/restore
func (h *FoodItemHandler) Restore(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromContext(r.Context())

	menuID := chi.URLParam(r, "menuID")
	menuIDInt, err := strconv.Atoi(menuID)
	if err != nil {
		http.Error(w, "invalid menu id", http.StatusBadRequest)
		return
	}
	if err := h.verifyMenuOwnership(r, int32(menuIDInt), claims.RestaurantID); err != nil {
		http.NotFound(w, r)
		return
	}

	itemID, err := strconv.Atoi(chi.URLParam(r, "itemID"))
	if err != nil {
		http.Error(w, "invalid item id", http.StatusBadRequest)
		return
	}
	if err := h.foodItemService.RestoreFoodItem(r.Context(), int32(itemID)); err != nil {
		http.Error(w, "failed to restore food item", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/manager/menus/"+menuID, http.StatusFound)
}

// AddIngredient handles POST /manager/menus/{menuID}/items/{itemID}/ingredients/add
func (h *FoodItemHandler) AddIngredient(w http.ResponseWriter, r *http.Request) {
	menuID := chi.URLParam(r, "menuID")
	itemID, err := strconv.Atoi(chi.URLParam(r, "itemID"))
	if err != nil {
		http.Error(w, "invalid item id", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form submission", http.StatusBadRequest)
		return
	}

	ingredientID, err := strconv.Atoi(r.FormValue("ingredient_id"))
	if err != nil {
		http.Redirect(w, r, "/manager/menus/"+menuID+"/items/"+strconv.Itoa(itemID)+"/edit", http.StatusFound)
		return
	}
	amount, _ := strconv.ParseFloat(r.FormValue("amount"), 64)
	measurement := strPtrOrNil(r.FormValue("measurement"))

	_ = h.ingredientItemService.AddIngredientToItem(r.Context(), int32(ingredientID), int32(itemID), measurement, &amount)

	http.Redirect(w, r, "/manager/menus/"+menuID+"/items/"+strconv.Itoa(itemID)+"/edit", http.StatusFound)
}

// CreateIngredient handles POST /manager/menus/{menuID}/items/{itemID}/ingredients/create
//
// This is the "search first, create if missing" flow: creates a brand
// new ingredient in the shared catalog, then immediately attaches it to
// this item — a manager can do both in one step rather than needing
// admin to add the ingredient first.
func (h *FoodItemHandler) CreateIngredient(w http.ResponseWriter, r *http.Request) {
	menuID := chi.URLParam(r, "menuID")
	itemID, err := strconv.Atoi(chi.URLParam(r, "itemID"))
	if err != nil {
		http.Error(w, "invalid item id", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form submission", http.StatusBadRequest)
		return
	}

	newName := r.FormValue("new_ingredient_name")
	if newName == "" {
		http.Redirect(w, r, "/manager/menus/"+menuID+"/items/"+strconv.Itoa(itemID)+"/edit", http.StatusFound)
		return
	}

	newIngredient, err := h.foodIngredientService.InsertFoodIngredient(r.Context(), newName, nil)
	if err != nil {
		http.Error(w, "failed to create ingredient", http.StatusInternalServerError)
		return
	}

	amount, _ := strconv.ParseFloat(r.FormValue("amount"), 64)
	measurement := strPtrOrNil(r.FormValue("measurement"))

	_ = h.ingredientItemService.AddIngredientToItem(r.Context(), newIngredient.ID, int32(itemID), measurement, &amount)

	http.Redirect(w, r, "/manager/menus/"+menuID+"/items/"+strconv.Itoa(itemID)+"/edit", http.StatusFound)
}

// RemoveIngredient handles POST /manager/menus/{menuID}/items/{itemID}/ingredients/{ingredientID}/remove
func (h *FoodItemHandler) RemoveIngredient(w http.ResponseWriter, r *http.Request) {
	menuID := chi.URLParam(r, "menuID")
	itemID, err := strconv.Atoi(chi.URLParam(r, "itemID"))
	if err != nil {
		http.Error(w, "invalid item id", http.StatusBadRequest)
		return
	}
	ingredientID, err := strconv.Atoi(chi.URLParam(r, "ingredientID"))
	if err != nil {
		http.Error(w, "invalid ingredient id", http.StatusBadRequest)
		return
	}

	_ = h.ingredientItemService.RemoveIngredientFromItem(r.Context(), int32(ingredientID), int32(itemID))

	http.Redirect(w, r, "/manager/menus/"+menuID+"/items/"+strconv.Itoa(itemID)+"/edit", http.StatusFound)
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

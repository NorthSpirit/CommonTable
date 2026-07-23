package admin

import (
	"net/http"

	"CommonTable/internal/services/foodorderservices"
	"CommonTable/internal/services/foodservices"
	"CommonTable/internal/services/restaurantservices"
	"CommonTable/internal/templates/admin/adminshared"
	"CommonTable/internal/templates/admin/dashboard"
)

type DashboardHandler struct {
	restaurantService     *restaurantservices.RestaurantService
	foodMenuService       *foodservices.FoodMenuService
	foodItemService       *foodservices.FoodItemService
	employeeService       *restaurantservices.EmployeeService
	foodIngredientService *foodservices.FoodIngredientService
	foodAllergenService   *foodservices.FoodAllergenService
	tableService          *foodorderservices.TableService
}

func NewDashboardHandler(
	restaurantService *restaurantservices.RestaurantService,
	foodMenuService *foodservices.FoodMenuService,
	foodItemService *foodservices.FoodItemService,
	employeeService *restaurantservices.EmployeeService,
	foodIngredientService *foodservices.FoodIngredientService,
	foodAllergenService *foodservices.FoodAllergenService,
	tableService *foodorderservices.TableService,
) *DashboardHandler {
	return &DashboardHandler{
		restaurantService:     restaurantService,
		foodMenuService:       foodMenuService,
		foodItemService:       foodItemService,
		employeeService:       employeeService,
		foodIngredientService: foodIngredientService,
		foodAllergenService:   foodAllergenService,
		tableService:          tableService,
	}
}

func (h *DashboardHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	restaurants, err := h.restaurantService.GetRestaurantList(r.Context())
	if err != nil {
		http.Error(w, "failed to load dashboard", http.StatusInternalServerError)
		return
	}

	menuCount, err := h.foodMenuService.CountFoodMenus(r.Context())
	if err != nil {
		http.Error(w, "failed to load dashboard", http.StatusInternalServerError)
		return
	}

	itemCount, err := h.foodItemService.CountFoodItems(r.Context())
	if err != nil {
		http.Error(w, "failed to load dashboard", http.StatusInternalServerError)
		return
	}

	employeeCount, err := h.employeeService.CountEmployees(r.Context())
	if err != nil {
		http.Error(w, "failed to load dashboard", http.StatusInternalServerError)
		return
	}

	ingredients, err := h.foodIngredientService.GetFoodIngredientList(r.Context())
	if err != nil {
		http.Error(w, "failed to load dashboard", http.StatusInternalServerError)
		return
	}

	allergens, err := h.foodAllergenService.GetFoodAllergensList(r.Context())
	if err != nil {
		http.Error(w, "failed to load dashboard", http.StatusInternalServerError)
		return
	}

	tables, err := h.tableService.ListTables(r.Context())
	if err != nil {
		http.Error(w, "failed to load dashboard", http.StatusInternalServerError)
		return
	}

	stats := adminshared.DashboardStats{
		Restaurants: int64(len(restaurants)),
		Menus:       menuCount,
		FoodItems:   itemCount,
		Employees:   employeeCount,
		Ingredients: int64(len(ingredients)),
		Allergens:   int64(len(allergens)),
		Tables:      int64(len(tables)),
	}

	if err := dashboard.DashboardPage(stats).Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

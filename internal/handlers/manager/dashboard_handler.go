package manager

import (
	"net/http"

	"CommonTable/internal/repository"
	"CommonTable/internal/services/foodorderservices"
	"CommonTable/internal/services/foodservices"
	"CommonTable/internal/services/restaurantservices"
	"CommonTable/internal/templates/manager"
)

type DashboardHandler struct {
	restaurantService *restaurantservices.RestaurantService
	foodMenuService   *foodservices.FoodMenuService
	foodItemService   *foodservices.FoodItemService
	employeeService   *restaurantservices.EmployeeService
	foodOrderService  *foodorderservices.FoodOrderService
}

func NewDashboardHandler(
	restaurantService *restaurantservices.RestaurantService,
	foodMenuService *foodservices.FoodMenuService,
	foodItemService *foodservices.FoodItemService,
	employeeService *restaurantservices.EmployeeService,
	foodOrderService *foodorderservices.FoodOrderService,
) *DashboardHandler {
	return &DashboardHandler{
		restaurantService: restaurantService,
		foodMenuService:   foodMenuService,
		foodItemService:   foodItemService,
		employeeService:   employeeService,
		foodOrderService:  foodOrderService,
	}
}

// Dashboard handles GET /manager
func (h *DashboardHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
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

	items, err := h.foodItemService.ListFoodItemsByRestaurantForAdmin(r.Context(), claims.RestaurantID)
	if err != nil {
		http.Error(w, "failed to load items", http.StatusInternalServerError)
		return
	}
	activeItems := int64(0)
	for _, item := range items {
		if !item.DeletedAt.Valid {
			activeItems++
		}
	}

	employees, err := h.employeeService.ListEmployeesByRestaurant(r.Context(), claims.RestaurantID)
	if err != nil {
		http.Error(w, "failed to load employees", http.StatusInternalServerError)
		return
	}
	activeStaff := int64(0)
	for _, e := range employees {
		if string(e.Role) != "manager" && e.IsActive {
			activeStaff++
		}
	}

	orders, err := h.foodOrderService.ListOrdersByRestaurantForAdmin(r.Context(), claims.RestaurantID)
	if err != nil {
		http.Error(w, "failed to load orders", http.StatusInternalServerError)
		return
	}
	activeOrders := int64(0)
	for _, o := range orders {
		if o.OrderStatus != repository.OrderStatusDelivered && o.OrderStatus != repository.OrderStatusCancelled {
			activeOrders++
		}
	}

	stats := manager.DashboardStats{
		Menus:     int64(len(menus)),
		FoodItems: activeItems,
		Employees: activeStaff,
		Orders:    activeOrders,
	}

	if err := manager.DashboardPage(restaurant.Name, claims.Name, stats).Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

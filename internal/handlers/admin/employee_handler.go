package admin

import (
	"errors"
	"net/http"
	"strconv"

	"CommonTable/internal/helpers"
	"CommonTable/internal/repository"
	"CommonTable/internal/services/restaurantservices"
	"CommonTable/internal/templates/admin/adminshared"
	"CommonTable/internal/templates/admin/employee"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

type EmployeeHandler struct {
	employeeService   *restaurantservices.EmployeeService
	restaurantService *restaurantservices.RestaurantService
}

func NewEmployeeHandler(employeeService *restaurantservices.EmployeeService, restaurantService *restaurantservices.RestaurantService) *EmployeeHandler {
	return &EmployeeHandler{employeeService: employeeService, restaurantService: restaurantService}
}

// List handles GET /admin/restaurants/{id}/employees
func (h *EmployeeHandler) List(w http.ResponseWriter, r *http.Request) {
	restaurantID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid restaurant id", http.StatusBadRequest)
		return
	}

	dbRestaurant, err := h.restaurantService.GetRestaurantByID(r.Context(), int32(restaurantID))
	if err != nil {
		http.Error(w, "failed to load restaurant", http.StatusInternalServerError)
		return
	}

	employees, err := h.employeeService.ListEmployeesByRestaurant(r.Context(), int32(restaurantID))
	if err != nil {
		http.Error(w, "failed to load employees", http.StatusInternalServerError)
		return
	}

	views := make([]adminshared.EmployeeView, 0, len(employees))
	for _, e := range employees {
		views = append(views, toEmployeeView(e))
	}

	if err := employee.EmployeesPage(int32(restaurantID), dbRestaurant.Name, views).Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

// NewForm handles GET /admin/restaurants/{id}/employees/new
func (h *EmployeeHandler) NewForm(w http.ResponseWriter, r *http.Request) {
	restaurantID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid restaurant id", http.StatusBadRequest)
		return
	}
	if err := employee.EmployeeFormPage(int32(restaurantID), false, adminshared.EmployeeView{}, "").Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

// Create handles POST /admin/restaurants/{id}/employees/new
func (h *EmployeeHandler) Create(w http.ResponseWriter, r *http.Request) {
	restaurantID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid restaurant id", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form submission", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	pin := r.FormValue("pin")
	roleRaw := r.FormValue("role")

	if name == "" || pin == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = employee.EmployeeFormPage(int32(restaurantID), false, adminshared.EmployeeView{Name: name, Role: roleRaw}, "Name and PIN are required.").Render(r.Context(), w)
		return
	}

	pinHash, err := helpers.HashPassword(pin)
	if err != nil {
		http.Error(w, "failed to hash pin", http.StatusInternalServerError)
		return
	}

	_, err = h.employeeService.InsertEmployee(
		r.Context(),
		int32(restaurantID),
		name,
		repository.EmployeeRole(roleRaw),
		&pinHash,
		nil,
		strPtrOrNil(r.FormValue("email")),
		strPtrOrNil(r.FormValue("phone_number")),
	)
	if err != nil {
		http.Error(w, "failed to create employee", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/restaurants/"+strconv.Itoa(restaurantID)+"/employees", http.StatusFound)
}

// EditForm handles GET /admin/restaurants/{id}/employees/{employeeID}/edit
func (h *EmployeeHandler) EditForm(w http.ResponseWriter, r *http.Request) {
	restaurantID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid restaurant id", http.StatusBadRequest)
		return
	}
	employeeID, err := strconv.Atoi(chi.URLParam(r, "employeeID"))
	if err != nil {
		http.Error(w, "invalid employee id", http.StatusBadRequest)
		return
	}

	dbEmployee, err := h.employeeService.GetEmployeeByID(r.Context(), int32(employeeID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "failed to load employee", http.StatusInternalServerError)
		return
	}

	if err := employee.EmployeeFormPage(int32(restaurantID), true, toEmployeeView(dbEmployee), "").Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

// Update handles POST /admin/restaurants/{id}/employees/{employeeID}/edit
func (h *EmployeeHandler) Update(w http.ResponseWriter, r *http.Request) {
	restaurantID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid restaurant id", http.StatusBadRequest)
		return
	}
	employeeID, err := strconv.Atoi(chi.URLParam(r, "employeeID"))
	if err != nil {
		http.Error(w, "invalid employee id", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form submission", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	roleRaw := r.FormValue("role")
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = employee.EmployeeFormPage(int32(restaurantID), true, adminshared.EmployeeView{ID: int32(employeeID), Name: name, Role: roleRaw}, "Name is required.").Render(r.Context(), w)
		return
	}

	err = h.employeeService.UpdateEmployee(
		r.Context(),
		int32(employeeID),
		name,
		repository.EmployeeRole(roleRaw),
		strPtrOrNil(r.FormValue("email")),
		strPtrOrNil(r.FormValue("phone_number")),
	)
	if err != nil {
		http.Error(w, "failed to update employee", http.StatusInternalServerError)
		return
	}

	if pin := r.FormValue("pin"); pin != "" {
		pinHash, err := helpers.HashPassword(pin)
		if err != nil {
			http.Error(w, "failed to hash pin", http.StatusInternalServerError)
			return
		}
		if err := h.employeeService.UpdateEmployeePinHash(r.Context(), int32(employeeID), pinHash); err != nil {
			http.Error(w, "failed to update pin", http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, "/admin/restaurants/"+strconv.Itoa(restaurantID)+"/employees", http.StatusFound)
}

// ToggleActive handles POST /admin/restaurants/{id}/employees/{employeeID}/toggle-active
func (h *EmployeeHandler) ToggleActive(w http.ResponseWriter, r *http.Request) {
	restaurantID := chi.URLParam(r, "id")
	employeeID, err := strconv.Atoi(chi.URLParam(r, "employeeID"))
	if err != nil {
		http.Error(w, "invalid employee id", http.StatusBadRequest)
		return
	}

	dbEmployee, err := h.employeeService.GetEmployeeByID(r.Context(), int32(employeeID))
	if err != nil {
		http.Error(w, "failed to load employee", http.StatusInternalServerError)
		return
	}

	if err := h.employeeService.SetEmployeeActive(r.Context(), int32(employeeID), !dbEmployee.IsActive); err != nil {
		http.Error(w, "failed to update employee status", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/restaurants/"+restaurantID+"/employees", http.StatusFound)
}

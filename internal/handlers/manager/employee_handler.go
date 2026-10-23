package manager

import (
	"errors"
	"net/http"
	"strconv"

	"CommonTable/internal/helpers"
	"CommonTable/internal/repository"
	"CommonTable/internal/services/restaurantservices"
	managertmpl "CommonTable/internal/templates/manager"

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

// List handles GET /manager/employees
// List handles GET /manager/employees
func (h *EmployeeHandler) List(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromContext(r.Context())

	restaurant, err := h.restaurantService.GetRestaurantByID(r.Context(), claims.RestaurantID)
	if err != nil {
		http.Error(w, "failed to load restaurant", http.StatusInternalServerError)
		return
	}

	employees, err := h.employeeService.ListEmployeesByRestaurant(r.Context(), claims.RestaurantID)
	if err != nil {
		http.Error(w, "failed to load employees", http.StatusInternalServerError)
		return
	}

	// Managers only manage waiter/chef staff — other managers (including
	// themselves) are never shown here, so there's no button that could
	// even be clicked to try deactivating a manager.
	views := make([]managertmpl.EmployeeView, 0, len(employees))
	for _, e := range employees {
		if string(e.Role) == "manager" {
			continue
		}
		views = append(views, toEmployeeView(e))
	}

	if err := managertmpl.EmployeesPage(restaurant.Name, claims.Name, views).Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

// NewForm handles GET /manager/employees/new
func (h *EmployeeHandler) NewForm(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromContext(r.Context())
	if err := managertmpl.EmployeeFormPage(claims.Name, false, managertmpl.EmployeeView{}, "").Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

// Create handles POST /manager/employees/new
func (h *EmployeeHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromContext(r.Context())

	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form submission", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	pin := r.FormValue("pin")
	roleRaw := r.FormValue("role")
	email := r.FormValue("email")
	phone := r.FormValue("phone_number")
	view := managertmpl.EmployeeView{Name: name, Role: roleRaw, Email: email, PhoneNumber: phone}

	// Enforced here too, not just hidden from the UI — a manager can
	// only ever create waiter or chef accounts, never another manager.
	if roleRaw != "waiter" && roleRaw != "chef" {
		w.WriteHeader(http.StatusBadRequest)
		_ = managertmpl.EmployeeFormPage(claims.Name, false, view, "Managers can only create waiter or chef accounts.").Render(r.Context(), w)
		return
	}

	if name == "" || pin == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = managertmpl.EmployeeFormPage(claims.Name, false, view, "Name and PIN are required.").Render(r.Context(), w)
		return
	}

	msg := helpers.CheckLengths(helpers.FieldLimit{Label: "Name", Value: name, Max: 255})
	if msg == "" {
		msg = helpers.ValidateContact(email, phone)
	}
	if msg == "" {
		msg = helpers.ValidateSecret("PIN", pin)
	}
	if msg != "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = managertmpl.EmployeeFormPage(claims.Name, false, view, msg).Render(r.Context(), w)
		return
	}

	pinHash, err := helpers.HashPassword(pin)
	if err != nil {
		http.Error(w, "failed to hash pin", http.StatusInternalServerError)
		return
	}

	_, err = h.employeeService.InsertEmployee(
		r.Context(),
		claims.RestaurantID,
		name,
		repository.EmployeeRole(roleRaw),
		&pinHash,
		nil,
		strPtrOrNil(email),
		strPtrOrNil(phone),
	)
	if err != nil {
		if isUniqueViolation(err) {
			w.WriteHeader(http.StatusBadRequest)
			_ = managertmpl.EmployeeFormPage(claims.Name, false, view, "An active employee named \""+name+"\" already exists at this restaurant.").Render(r.Context(), w)
			return
		}
		if msg := helpers.PgUserMessage(err); msg != "" {
			w.WriteHeader(http.StatusBadRequest)
			_ = managertmpl.EmployeeFormPage(claims.Name, false, view, msg).Render(r.Context(), w)
			return
		}
		http.Error(w, "failed to create employee", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/manager/employees", http.StatusFound)
}

// EditForm handles GET /manager/employees/{employeeID}/edit
func (h *EmployeeHandler) EditForm(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromContext(r.Context())

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

	// A manager may only edit employees belonging to their own
	// restaurant — without this check, a manager could edit any
	// employee ID at any restaurant just by guessing/incrementing URLs.
	if dbEmployee.RestaurantID != claims.RestaurantID {
		http.NotFound(w, r)
		return
	}

	// Managers cannot edit other managers.
	if string(dbEmployee.Role) == "manager" {
		http.NotFound(w, r)
		return
	}

	if err := managertmpl.EmployeeFormPage(claims.Name, true, toEmployeeView(dbEmployee), "").Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

// Update handles POST /manager/employees/{employeeID}/edit
func (h *EmployeeHandler) Update(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromContext(r.Context())

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
	if dbEmployee.RestaurantID != claims.RestaurantID || string(dbEmployee.Role) == "manager" {
		http.NotFound(w, r)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form submission", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	roleRaw := r.FormValue("role")
	email := r.FormValue("email")
	phone := r.FormValue("phone_number")
	pin := r.FormValue("pin")
	view := managertmpl.EmployeeView{ID: int32(employeeID), Name: name, Role: roleRaw, Email: email, PhoneNumber: phone}
	if roleRaw != "waiter" && roleRaw != "chef" {
		w.WriteHeader(http.StatusBadRequest)
		_ = managertmpl.EmployeeFormPage(claims.Name, true, view, "Managers can only assign waiter or chef roles.").Render(r.Context(), w)
		return
	}
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = managertmpl.EmployeeFormPage(claims.Name, true, view, "Name is required.").Render(r.Context(), w)
		return
	}

	msg := helpers.CheckLengths(helpers.FieldLimit{Label: "Name", Value: name, Max: 255})
	if msg == "" {
		msg = helpers.ValidateContact(email, phone)
	}
	if msg == "" {
		msg = helpers.ValidateSecret("PIN", pin)
	}
	if msg != "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = managertmpl.EmployeeFormPage(claims.Name, true, view, msg).Render(r.Context(), w)
		return
	}

	err = h.employeeService.UpdateEmployee(
		r.Context(),
		int32(employeeID),
		name,
		repository.EmployeeRole(roleRaw),
		strPtrOrNil(email),
		strPtrOrNil(phone),
	)
	if err != nil {
		if isUniqueViolation(err) {
			w.WriteHeader(http.StatusBadRequest)
			_ = managertmpl.EmployeeFormPage(claims.Name, true, view, "An active employee named \""+name+"\" already exists at this restaurant.").Render(r.Context(), w)
			return
		}
		if msg := helpers.PgUserMessage(err); msg != "" {
			w.WriteHeader(http.StatusBadRequest)
			_ = managertmpl.EmployeeFormPage(claims.Name, true, view, msg).Render(r.Context(), w)
			return
		}
		http.Error(w, "failed to update employee", http.StatusInternalServerError)
		return
	}

	if pin != "" {
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

	http.Redirect(w, r, "/manager/employees", http.StatusFound)
}

// ToggleActive handles POST /manager/employees/{employeeID}/toggle-active
func (h *EmployeeHandler) ToggleActive(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromContext(r.Context())

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
	if dbEmployee.RestaurantID != claims.RestaurantID || string(dbEmployee.Role) == "manager" {
		http.NotFound(w, r)
		return
	}

	if err := h.employeeService.SetEmployeeActive(r.Context(), int32(employeeID), !dbEmployee.IsActive); err != nil {
		http.Error(w, "failed to update employee status", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/manager/employees", http.StatusFound)
}

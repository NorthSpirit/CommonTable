package staff

import (
	"net/http"

	"CommonTable/internal/helpers"
	"CommonTable/internal/services/restaurantservices"
	stafftmpl "CommonTable/internal/templates/staff"
)

type ProfileHandler struct {
	employeeService *restaurantservices.EmployeeService
}

func NewProfileHandler(employeeService *restaurantservices.EmployeeService) *ProfileHandler {
	return &ProfileHandler{employeeService: employeeService}
}

// Edit handles GET /staff/profile
func (h *ProfileHandler) Edit(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromContext(r.Context())

	dbEmployee, err := h.employeeService.GetEmployeeByID(r.Context(), claims.EmployeeID)
	if err != nil {
		http.Error(w, "failed to load profile", http.StatusInternalServerError)
		return
	}

	email := ""
	if dbEmployee.Email.Valid {
		email = dbEmployee.Email.String
	}
	phone := ""
	if dbEmployee.PhoneNumber.Valid {
		phone = dbEmployee.PhoneNumber.String
	}

	if err := stafftmpl.ProfilePage(claims.Name, email, phone, "", "").Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

// Update handles POST /staff/profile
func (h *ProfileHandler) Update(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromContext(r.Context())

	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form submission", http.StatusBadRequest)
		return
	}

	dbEmployee, err := h.employeeService.GetEmployeeByID(r.Context(), claims.EmployeeID)
	if err != nil {
		http.Error(w, "failed to load profile", http.StatusInternalServerError)
		return
	}

	email := r.FormValue("email")
	phone := r.FormValue("phone_number")

	if err := h.employeeService.UpdateEmployee(r.Context(), claims.EmployeeID, dbEmployee.Name, dbEmployee.Role, strPtrOrNil(email), strPtrOrNil(phone)); err != nil {
		http.Error(w, "failed to update profile", http.StatusInternalServerError)
		return
	}

	if newPin := r.FormValue("new_pin"); newPin != "" {
		pinHash, err := helpers.HashPassword(newPin)
		if err != nil {
			http.Error(w, "failed to hash pin", http.StatusInternalServerError)
			return
		}
		if err := h.employeeService.UpdateEmployeePinHash(r.Context(), claims.EmployeeID, pinHash); err != nil {
			http.Error(w, "failed to update pin", http.StatusInternalServerError)
			return
		}
	}

	if err := stafftmpl.ProfilePage(claims.Name, email, phone, "Profile updated.", "").Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

func strPtrOrNil(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

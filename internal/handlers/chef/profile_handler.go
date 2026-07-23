package chef

import (
	"net/http"

	"CommonTable/internal/helpers"
	"CommonTable/internal/services/restaurantservices"
	"CommonTable/internal/templates/chef"
)

type ProfileHandler struct {
	employeeService *restaurantservices.EmployeeService
}

func NewProfileHandler(employeeService *restaurantservices.EmployeeService) *ProfileHandler {
	return &ProfileHandler{employeeService: employeeService}
}

// Edit handles GET /chef/profile
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

	if err := chef.ProfilePage(claims.Name, email, phone, "", "").Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

// Update handles POST /chef/profile
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

	// Name and role are deliberately not editable here — only an admin
	// changes those, via the employee's admin record.
	err = h.employeeService.UpdateEmployee(
		r.Context(),
		claims.EmployeeID,
		dbEmployee.Name,
		dbEmployee.Role,
		strPtrOrNil(email),
		strPtrOrNil(phone),
	)
	if err != nil {
		http.Error(w, "failed to update profile", http.StatusInternalServerError)
		return
	}

	newPin := r.FormValue("new_pin")
	newPinConfirm := r.FormValue("new_pin_confirm")

	if newPin != "" || newPinConfirm != "" {
		if newPin != newPinConfirm {
			w.WriteHeader(http.StatusBadRequest)
			_ = chef.ProfilePage(claims.Name, email, phone, "", "New PIN and confirmation do not match.").Render(r.Context(), w)
			return
		}

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

	if err := chef.ProfilePage(claims.Name, email, phone, "Profile updated.", "").Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

func strPtrOrNil(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

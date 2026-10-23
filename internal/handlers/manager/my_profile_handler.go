package manager

import (
	"net/http"

	"CommonTable/internal/helpers"
	"CommonTable/internal/services/restaurantservices"
	managertmpl "CommonTable/internal/templates/manager"
)

type MyProfileHandler struct {
	employeeService *restaurantservices.EmployeeService
}

func NewMyProfileHandler(employeeService *restaurantservices.EmployeeService) *MyProfileHandler {
	return &MyProfileHandler{employeeService: employeeService}
}

// Edit handles GET /manager/my-profile
func (h *MyProfileHandler) Edit(w http.ResponseWriter, r *http.Request) {
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

	if err := managertmpl.MyProfilePage(claims.Name, email, phone, "", "").Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

// Update handles POST /manager/my-profile
func (h *MyProfileHandler) Update(w http.ResponseWriter, r *http.Request) {
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

	if msg := helpers.ValidateContact(email, phone); msg != "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = managertmpl.MyProfilePage(claims.Name, email, phone, "", msg).Render(r.Context(), w)
		return
	}

	newPin := r.FormValue("new_pin")
	if len(newPin) > 72 {
		w.WriteHeader(http.StatusBadRequest)
		_ = managertmpl.MyProfilePage(claims.Name, email, phone, "", "PIN must be at most 72 characters.").Render(r.Context(), w)
		return
	}
	
	if err := h.employeeService.UpdateEmployee(r.Context(), claims.EmployeeID, dbEmployee.Name, dbEmployee.Role, strPtrOrNil(email), strPtrOrNil(phone)); err != nil {
		http.Error(w, "failed to update profile", http.StatusInternalServerError)
		return
	}

	if newPin != "" {
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

	if err := managertmpl.MyProfilePage(claims.Name, email, phone, "Profile updated.", "").Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

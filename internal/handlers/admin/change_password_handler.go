package admin

import (
	"errors"
	"net/http"

	"CommonTable/internal/services/adminservices"
	"CommonTable/internal/templates/admin/adminshared"
)

type ChangePasswordHandler struct {
	adminService *adminservices.AdminService
}

func NewChangePasswordHandler(adminService *adminservices.AdminService) *ChangePasswordHandler {
	return &ChangePasswordHandler{adminService: adminService}
}

// Edit handles GET /admin/change-password
func (h *ChangePasswordHandler) Edit(w http.ResponseWriter, r *http.Request) {
	if err := adminshared.ChangePasswordPage("", "").Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

// Update handles POST /admin/change-password
func (h *ChangePasswordHandler) Update(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form submission", http.StatusBadRequest)
		return
	}

	currentPassword := r.FormValue("current_password")
	newPassword := r.FormValue("new_password")
	newPasswordConfirm := r.FormValue("new_password_confirm")

	if newPassword != newPasswordConfirm {
		w.WriteHeader(http.StatusBadRequest)
		_ = adminshared.ChangePasswordPage("", "New password and confirmation do not match.").Render(r.Context(), w)
		return
	}

	if len(newPassword) < 8 {
		w.WriteHeader(http.StatusBadRequest)
		_ = adminshared.ChangePasswordPage("", "New password must be at least 8 characters.").Render(r.Context(), w)
		return
	}

	err := h.adminService.ChangePassword(currentPassword, newPassword)
	if err != nil {
		if errors.Is(err, adminservices.ErrInvalidCredentials) {
			w.WriteHeader(http.StatusBadRequest)
			_ = adminshared.ChangePasswordPage("", "Current password is incorrect.").Render(r.Context(), w)
			return
		}
		http.Error(w, "failed to change password", http.StatusInternalServerError)
		return
	}

	if err := adminshared.ChangePasswordPage("Password changed successfully.", "").Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

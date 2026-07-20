package admin

import (
	"errors"
	"net/http"
	"time"

	"CommonTable/internal/services/adminservices"
	"CommonTable/internal/templates/admin"
)

type AuthHandler struct {
	adminService *adminservices.AdminService
	tokenTTL     time.Duration
}

func NewAuthHandler(adminService *adminservices.AdminService, tokenTTL time.Duration) *AuthHandler {
	return &AuthHandler{adminService: adminService, tokenTTL: tokenTTL}
}

// LoginPage handles GET /admin/login
func (h *AuthHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	if err := admin.LoginPage("").Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

// Login handles POST /admin/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form submission", http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	err := h.adminService.VerifyAdminLogin(username, password)
	if err != nil {
		if errors.Is(err, adminservices.ErrInvalidCredentials) {
			_ = admin.LoginPage("Invalid username or password.").Render(r.Context(), w)
			return
		}
		http.Error(w, "login failed", http.StatusInternalServerError)
		return
	}

	token, err := h.adminService.IssueAdminToken()
	if err != nil {
		http.Error(w, "failed to issue session", http.StatusInternalServerError)
		return
	}

	setSessionCookie(w, token, h.tokenTTL)
	http.Redirect(w, r, "/admin", http.StatusFound)
}

// Logout handles GET /admin/logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	clearSessionCookie(w)
	http.Redirect(w, r, "/admin/login", http.StatusFound)
}

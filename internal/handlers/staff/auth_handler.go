package staff

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"CommonTable/internal/services/authservices"
	"CommonTable/internal/services/restaurantservices"
	stafftmpl "CommonTable/internal/templates/staff"

	"github.com/jackc/pgx/v5"
)

type AuthHandler struct {
	authService       *authservices.AuthService
	restaurantService *restaurantservices.RestaurantService
	tokenTTL          time.Duration
}

func NewAuthHandler(authService *authservices.AuthService, restaurantService *restaurantservices.RestaurantService, tokenTTL time.Duration) *AuthHandler {
	return &AuthHandler{authService: authService, restaurantService: restaurantService, tokenTTL: tokenTTL}
}

// RestaurantLoginPage handles GET /staff/login
func (h *AuthHandler) RestaurantLoginPage(w http.ResponseWriter, r *http.Request) {
	if err := stafftmpl.RestaurantLoginPage("").Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

// RestaurantLogin handles POST /staff/login
func (h *AuthHandler) RestaurantLogin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form submission", http.StatusBadRequest)
		return
	}

	name := r.FormValue("restaurant_name")
	password := r.FormValue("password")

	restaurant, err := h.authService.VerifyRestaurantPassword(r.Context(), name, password)
	if err != nil {
		if errors.Is(err, authservices.ErrInvalidCredentials) || errors.Is(err, pgx.ErrNoRows) {
			w.WriteHeader(http.StatusBadRequest)
			_ = stafftmpl.RestaurantLoginPage("Incorrect restaurant name or password.").Render(r.Context(), w)
			return
		}
		http.Error(w, "login failed", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/staff/login/employee?restaurant="+strconv.Itoa(int(restaurant.ID)), http.StatusFound)
}

// EmployeeLoginPage handles GET /staff/login/employee
func (h *AuthHandler) EmployeeLoginPage(w http.ResponseWriter, r *http.Request) {
	restaurantID, err := strconv.Atoi(r.URL.Query().Get("restaurant"))
	if err != nil {
		http.Redirect(w, r, "/staff/login", http.StatusFound)
		return
	}

	restaurant, err := h.restaurantService.GetRestaurantByID(r.Context(), int32(restaurantID))
	if err != nil {
		http.Redirect(w, r, "/staff/login", http.StatusFound)
		return
	}

	if err := stafftmpl.EmployeeLoginPage(int32(restaurantID), restaurant.Name, "").Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

// EmployeeLogin handles POST /staff/login/employee
func (h *AuthHandler) EmployeeLogin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form submission", http.StatusBadRequest)
		return
	}

	restaurantID, err := strconv.Atoi(r.FormValue("restaurant_id"))
	if err != nil {
		http.Redirect(w, r, "/staff/login", http.StatusFound)
		return
	}

	employeeName := r.FormValue("employee_name")
	password := r.FormValue("password")

	restaurant, err := h.restaurantService.GetRestaurantByID(r.Context(), int32(restaurantID))
	if err != nil {
		http.Redirect(w, r, "/staff/login", http.StatusFound)
		return
	}

	dbEmployee, err := h.authService.VerifyEmployeeLogin(r.Context(), int32(restaurantID), employeeName, password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = stafftmpl.EmployeeLoginPage(int32(restaurantID), restaurant.Name, "Incorrect name or password.").Render(r.Context(), w)
		return
	}

	// This is the enforcement point for "a chef/manager can't log into
	// the waiter web page" — the credentials may be perfectly valid, but
	// only the waiter role is allowed through here. The message points
	// them at the right front door instead of just refusing.
	role := string(dbEmployee.Role)
	if role != "waiter" {
		w.WriteHeader(http.StatusForbidden)
		message := "This page is for waiters only."
		switch role {
		case "chef":
			message = "This login is for waiters. Chefs should sign in at /chef/login."
		case "manager":
			message = "This login is for waiters. Managers should sign in at /manager/login."
		}
		_ = stafftmpl.EmployeeLoginPage(int32(restaurantID), restaurant.Name, message).Render(r.Context(), w)
		return
	}

	token, err := h.authService.IssueEmployeeToken(dbEmployee)
	if err != nil {
		http.Error(w, "failed to issue session", http.StatusInternalServerError)
		return
	}

	setSessionCookie(w, token, h.tokenTTL)
	http.Redirect(w, r, "/staff/profile", http.StatusFound)
}

// Logout handles GET /staff/logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	clearSessionCookie(w)
	http.Redirect(w, r, "/staff/login", http.StatusFound)
}

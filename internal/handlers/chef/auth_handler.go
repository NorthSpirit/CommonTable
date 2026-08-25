package chef

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"CommonTable/internal/services/authservices"
	"CommonTable/internal/services/restaurantservices"
	"CommonTable/internal/templates/chef"

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

// RestaurantLoginPage handles GET /chef/login
func (h *AuthHandler) RestaurantLoginPage(w http.ResponseWriter, r *http.Request) {
	if err := chef.RestaurantLoginPage("").Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

// RestaurantLogin handles POST /chef/login
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
			_ = chef.RestaurantLoginPage("Incorrect restaurant name or password.").Render(r.Context(), w)
			return
		}
		http.Error(w, "login failed", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/chef/login/employee?restaurant="+strconv.Itoa(int(restaurant.ID)), http.StatusFound)
}

// EmployeeLoginPage handles GET /chef/login/employee
func (h *AuthHandler) EmployeeLoginPage(w http.ResponseWriter, r *http.Request) {
	restaurantID, err := strconv.Atoi(r.URL.Query().Get("restaurant"))
	if err != nil {
		http.Redirect(w, r, "/chef/login", http.StatusFound)
		return
	}

	restaurant, err := h.restaurantService.GetRestaurantByID(r.Context(), int32(restaurantID))
	if err != nil {
		http.Redirect(w, r, "/chef/login", http.StatusFound)
		return
	}

	if err := chef.EmployeeLoginPage(int32(restaurantID), restaurant.Name, "").Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

// EmployeeLogin handles POST /chef/login/employee
func (h *AuthHandler) EmployeeLogin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form submission", http.StatusBadRequest)
		return
	}

	restaurantID, err := strconv.Atoi(r.FormValue("restaurant_id"))
	if err != nil {
		http.Redirect(w, r, "/chef/login", http.StatusFound)
		return
	}

	employeeName := r.FormValue("employee_name")
	password := r.FormValue("password")

	restaurant, err := h.restaurantService.GetRestaurantByID(r.Context(), int32(restaurantID))
	if err != nil {
		http.Redirect(w, r, "/chef/login", http.StatusFound)
		return
	}

	dbEmployee, err := h.authService.VerifyEmployeeLogin(r.Context(), int32(restaurantID), employeeName, password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = chef.EmployeeLoginPage(int32(restaurantID), restaurant.Name, "Incorrect name or password.").Render(r.Context(), w)
		return
	}

	// This is the enforcement point for "a waiter/manager can't log into
	// the chef web page" — the credentials may be perfectly valid, but
	// only the chef role is allowed through here. The message points
	// them at the right front door instead of just refusing.
	role := string(dbEmployee.Role)
	if role != "chef" {
		w.WriteHeader(http.StatusForbidden)
		message := "This web page is for kitchen staff only."
		switch role {
		case "manager":
			message = "This login is for chefs. Managers should sign in at /manager/login."
		case "waiter":
			message = "This login is for chefs. Waiters can use the mobile app, or /staff/login for profile & chat on the web."
		}
		_ = chef.EmployeeLoginPage(int32(restaurantID), restaurant.Name, message).Render(r.Context(), w)
		return
	}

	token, err := h.authService.IssueEmployeeToken(dbEmployee)
	if err != nil {
		http.Error(w, "failed to issue session", http.StatusInternalServerError)
		return
	}

	setSessionCookie(w, token, h.tokenTTL)
	http.Redirect(w, r, "/chef/menus", http.StatusFound)
}

// Logout handles GET /chef/logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	clearSessionCookie(w)
	http.Redirect(w, r, "/chef/login", http.StatusFound)
}

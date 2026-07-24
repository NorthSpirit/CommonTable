package kitchen

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"CommonTable/internal/services/authservices"
	"CommonTable/internal/services/restaurantservices"
	"CommonTable/internal/templates/kitchen"

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

func (h *AuthHandler) RestaurantLoginPage(w http.ResponseWriter, r *http.Request) {
	if err := kitchen.RestaurantLoginPage("").Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

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
			_ = kitchen.RestaurantLoginPage("Incorrect restaurant name or password.").Render(r.Context(), w)
			return
		}
		http.Error(w, "login failed", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/kitchen/login/employee?restaurant="+strconv.Itoa(int(restaurant.ID)), http.StatusFound)
}

func (h *AuthHandler) EmployeeLoginPage(w http.ResponseWriter, r *http.Request) {
	restaurantID, err := strconv.Atoi(r.URL.Query().Get("restaurant"))
	if err != nil {
		http.Redirect(w, r, "/kitchen/login", http.StatusFound)
		return
	}

	restaurant, err := h.restaurantService.GetRestaurantByID(r.Context(), int32(restaurantID))
	if err != nil {
		http.Redirect(w, r, "/kitchen/login", http.StatusFound)
		return
	}

	if err := kitchen.EmployeeLoginPage(int32(restaurantID), restaurant.Name, "").Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

func (h *AuthHandler) EmployeeLogin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form submission", http.StatusBadRequest)
		return
	}

	restaurantID, err := strconv.Atoi(r.FormValue("restaurant_id"))
	if err != nil {
		http.Redirect(w, r, "/kitchen/login", http.StatusFound)
		return
	}

	employeeName := r.FormValue("employee_name")
	password := r.FormValue("password")

	restaurant, err := h.restaurantService.GetRestaurantByID(r.Context(), int32(restaurantID))
	if err != nil {
		http.Redirect(w, r, "/kitchen/login", http.StatusFound)
		return
	}

	dbEmployee, err := h.authService.VerifyEmployeeLogin(r.Context(), int32(restaurantID), employeeName, password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = kitchen.EmployeeLoginPage(int32(restaurantID), restaurant.Name, "Incorrect name or password.").Render(r.Context(), w)
		return
	}

	role := string(dbEmployee.Role)
	if role != "chef" && role != "manager" {
		w.WriteHeader(http.StatusForbidden)
		_ = kitchen.EmployeeLoginPage(int32(restaurantID), restaurant.Name, "The kitchen board is for chefs and managers only.").Render(r.Context(), w)
		return
	}

	token, err := h.authService.IssueEmployeeToken(dbEmployee)
	if err != nil {
		http.Error(w, "failed to issue session", http.StatusInternalServerError)
		return
	}

	setSessionCookie(w, token, h.tokenTTL)
	http.Redirect(w, r, "/kitchen/orders", http.StatusFound)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	clearSessionCookie(w)
	http.Redirect(w, r, "/kitchen/login", http.StatusFound)
}

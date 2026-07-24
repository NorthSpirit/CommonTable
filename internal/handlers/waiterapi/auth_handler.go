package waiterapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"CommonTable/internal/services/authservices"
)

type AuthHandler struct {
	authService *authservices.AuthService
}

func NewAuthHandler(authService *authservices.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type restaurantLoginRequest struct {
	RestaurantName string `json:"restaurant_name"`
	Password       string `json:"password"`
}

type restaurantLoginResponse struct {
	RestaurantID   int32  `json:"restaurant_id"`
	RestaurantName string `json:"restaurant_name"`
}

// RestaurantLogin handles POST /api/waiter/login
func (h *AuthHandler) RestaurantLogin(w http.ResponseWriter, r *http.Request) {
	var req restaurantLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.RestaurantName == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "restaurant_name and password are required")
		return
	}

	restaurant, err := h.authService.VerifyRestaurantPassword(r.Context(), req.RestaurantName, req.Password)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "incorrect restaurant name or password")
		return
	}

	writeJSON(w, http.StatusOK, restaurantLoginResponse{
		RestaurantID:   restaurant.ID,
		RestaurantName: restaurant.Name,
	})
}

type employeeLoginRequest struct {
	RestaurantID int32  `json:"restaurant_id"`
	EmployeeName string `json:"employee_name"`
	Password     string `json:"password"`
}

type employeeLoginResponse struct {
	Token      string    `json:"token"`
	EmployeeID int32     `json:"employee_id"`
	Name       string    `json:"name"`
	Role       string    `json:"role"`
	ExpiresAt  time.Time `json:"expires_at"`
}

// EmployeeLogin handles POST /api/waiter/login/employee
func (h *AuthHandler) EmployeeLogin(w http.ResponseWriter, r *http.Request) {
	var req employeeLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.RestaurantID == 0 || req.EmployeeName == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "restaurant_id, employee_name and password are required")
		return
	}

	dbEmployee, err := h.authService.VerifyEmployeeLogin(r.Context(), req.RestaurantID, req.EmployeeName, req.Password)
	if err != nil {
		if errors.Is(err, authservices.ErrInvalidCredentials) {
			writeError(w, http.StatusUnauthorized, "incorrect name or password")
			return
		}
		writeError(w, http.StatusInternalServerError, "login failed")
		return
	}

	role := string(dbEmployee.Role)
	if role != "waiter" && role != "manager" {
		writeError(w, http.StatusForbidden, "the waiter app is for waiters and managers only")
		return
	}

	token, err := h.authService.IssueEmployeeTokenWithTTL(dbEmployee, tokenTTL)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to issue session")
		return
	}

	writeJSON(w, http.StatusOK, employeeLoginResponse{
		Token:      token,
		EmployeeID: dbEmployee.ID,
		Name:       dbEmployee.Name,
		Role:       role,
		ExpiresAt:  time.Now().Add(tokenTTL),
	})
}

package authservices

import (
	"context"
	"errors"
	"fmt"
	"time"

	"CommonTable/internal/helpers"
	"CommonTable/internal/repository"

	"github.com/golang-jwt/jwt/v5"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type EmployeeClaims struct {
	RestaurantID int32  `json:"restaurant_id"`
	EmployeeID   int32  `json:"employee_id"`
	Role         string `json:"role"`
	jwt.RegisteredClaims
}

type AuthService struct {
	queries   *repository.Queries
	jwtSecret []byte
	tokenTTL  time.Duration
}

func NewAuthService(queries *repository.Queries, jwtSecret []byte, tokenTTL time.Duration) *AuthService {
	return &AuthService{queries: queries, jwtSecret: jwtSecret, tokenTTL: tokenTTL}
}

// VerifyRestaurantPassword handles step 1 of your two-step login: identify
// and authenticate the restaurant. It does not issue a token — the client
// only needs the restaurant_id to proceed to step 2.
func (a *AuthService) VerifyRestaurantPassword(ctx context.Context, name string, password string) (repository.Restaurant, error) {
	restaurant, err := a.queries.GetRestaurantByName(ctx, name)
	if err != nil {
		return repository.Restaurant{}, fmt.Errorf("get restaurant by name: %w", err)
	}

	if !restaurant.PasswordHash.Valid || !helpers.CheckPasswordHash(password, restaurant.PasswordHash.String) {
		return repository.Restaurant{}, ErrInvalidCredentials
	}

	return restaurant, nil
}

// VerifyEmployeeLogin handles step 2: authenticate the employee within the
// already-established restaurant.
func (a *AuthService) VerifyEmployeeLogin(ctx context.Context, restaurantID int32, employeeName string, password string) (repository.Employee, error) {
	employee, err := a.queries.GetEmployeeByRestaurantAndName(ctx, repository.GetEmployeeByRestaurantAndNameParams{
		RestaurantID: restaurantID,
		Name:         helpers.TextFromString(&employeeName),
	})
	if err != nil {
		return repository.Employee{}, fmt.Errorf("get employee by restaurant and name: %w", err)
	}

	if !employee.PinHash.Valid || !helpers.CheckPasswordHash(password, employee.PinHash.String) {
		return repository.Employee{}, ErrInvalidCredentials
	}

	return employee, nil
}

// IssueEmployeeToken creates a signed JWT carrying restaurant and employee
// identity — issued after both login steps succeed.
func (a *AuthService) IssueEmployeeToken(employee repository.Employee) (string, error) {
	claims := EmployeeClaims{
		RestaurantID: employee.RestaurantID,
		EmployeeID:   employee.ID,
		Role:         string(employee.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(a.tokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(a.jwtSecret)
}

// ValidateEmployeeToken parses and verifies an incoming request's token —
// this is what your auth middleware will call on every protected request.
func (a *AuthService) ValidateEmployeeToken(tokenString string) (*EmployeeClaims, error) {
	claims := &EmployeeClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return a.jwtSecret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("parse employee token: %w", err)
	}

	if !token.Valid {
		return nil, ErrInvalidCredentials
	}

	return claims, nil
}

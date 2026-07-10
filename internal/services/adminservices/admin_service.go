package adminservices

import (
	"errors"
	"time"

	"CommonTable/internal/helpers"

	"github.com/golang-jwt/jwt/v5"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

// AdminConfig is loaded from your config file — Username plaintext,
// PasswordHash pre-hashed offline (never store the raw admin password,
// even in a config file only you control).
type AdminConfig struct {
	Username     string
	PasswordHash string
}

type AdminClaims struct {
	Username string `json:"username"`
	IsAdmin  bool   `json:"is_admin"`
	jwt.RegisteredClaims
}

type AdminService struct {
	config    AdminConfig
	jwtSecret []byte
	tokenTTL  time.Duration
}

func NewAdminService(config AdminConfig, jwtSecret []byte, tokenTTL time.Duration) *AdminService {
	return &AdminService{config: config, jwtSecret: jwtSecret, tokenTTL: tokenTTL}
}

func (a *AdminService) VerifyAdminLogin(username string, password string) error {
	if username != a.config.Username {
		return ErrInvalidCredentials
	}

	if !helpers.CheckPasswordHash(password, a.config.PasswordHash) {
		return ErrInvalidCredentials
	}

	return nil
}

func (a *AdminService) IssueAdminToken() (string, error) {
	claims := AdminClaims{
		Username: a.config.Username,
		IsAdmin:  true,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(a.tokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(a.jwtSecret)
}

func (a *AdminService) ValidateAdminToken(tokenString string) (*AdminClaims, error) {
	claims := &AdminClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return a.jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, ErrInvalidCredentials
	}

	return claims, nil
}

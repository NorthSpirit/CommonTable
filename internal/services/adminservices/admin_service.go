package adminservices

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"CommonTable/internal/helpers"

	"github.com/golang-jwt/jwt/v5"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type AdminConfig struct {
	Username     string
	PasswordHash string
}

type AdminClaims struct {
	Username string `json:"username"`
	IsAdmin  bool   `json:"is_admin"`
	jwt.RegisteredClaims
}

type storedCredentials struct {
	PasswordHash string `json:"password_hash"`
}

type AdminService struct {
	mu              sync.RWMutex
	config          AdminConfig
	credentialsPath string
	jwtSecret       []byte
	tokenTTL        time.Duration
}

// NewAdminService loads the admin credential from credentialsPath if that
// file exists (meaning the password has been changed at least once since
// deployment); otherwise it falls back to the hash baked into config at
// startup, which comes from .env — the "factory default."
func NewAdminService(config AdminConfig, credentialsPath string, jwtSecret []byte, tokenTTL time.Duration) *AdminService {
	s := &AdminService{
		config:          config,
		credentialsPath: credentialsPath,
		jwtSecret:       jwtSecret,
		tokenTTL:        tokenTTL,
	}

	if stored, err := s.loadStoredCredentials(); err == nil {
		s.config.PasswordHash = stored.PasswordHash
	}

	return s
}

func (a *AdminService) loadStoredCredentials() (storedCredentials, error) {
	data, err := os.ReadFile(a.credentialsPath)
	if err != nil {
		return storedCredentials{}, err
	}

	var stored storedCredentials
	if err := json.Unmarshal(data, &stored); err != nil {
		return storedCredentials{}, err
	}

	return stored, nil
}

func (a *AdminService) VerifyAdminLogin(username string, password string) error {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if username != a.config.Username {
		return ErrInvalidCredentials
	}

	if !helpers.CheckPasswordHash(password, a.config.PasswordHash) {
		return ErrInvalidCredentials
	}

	return nil
}

// VerifyPassword checks a password against the admin's current hash,
// without checking username — used for re-confirming identity for
// sensitive actions (e.g. deleting a restaurant).
func (a *AdminService) VerifyPassword(password string) error {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if !helpers.CheckPasswordHash(password, a.config.PasswordHash) {
		return ErrInvalidCredentials
	}

	return nil
}

// ChangePassword verifies the current password, then hashes and persists
// the new one to disk, and updates the in-memory hash immediately so the
// change takes effect without a server restart.
func (a *AdminService) ChangePassword(currentPassword string, newPassword string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if !helpers.CheckPasswordHash(currentPassword, a.config.PasswordHash) {
		return ErrInvalidCredentials
	}

	newHash, err := helpers.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("hash new password: %w", err)
	}

	stored := storedCredentials{PasswordHash: newHash}
	data, err := json.MarshalIndent(stored, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal credentials: %w", err)
	}

	if err := os.WriteFile(a.credentialsPath, data, 0600); err != nil {
		return fmt.Errorf("write credentials file: %w", err)
	}

	a.config.PasswordHash = newHash
	return nil
}

func (a *AdminService) IssueAdminToken() (string, error) {
	a.mu.RLock()
	username := a.config.Username
	a.mu.RUnlock()

	claims := AdminClaims{
		Username: username,
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

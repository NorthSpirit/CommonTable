package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DBString string

	AdminUsername        string
	AdminPasswordHash    string
	AdminCredentialsPath string

	JWTSecret   []byte
	JWTTokenTTL time.Duration

	AltLanguageLabel string
	TranslationsDir  string
	UploadsDir       string

	DeploymentMode string
	PublicBaseURL  string
}

func Load() (*Config, error) {
	// Fall back in case .env is flawed
	_ = godotenv.Load()

	dbString := os.Getenv("GOOSE_DBSTRING")
	if dbString == "" {
		return nil, fmt.Errorf("GOOSE_DBSTRING is not set")
	}

	adminUsername := os.Getenv("ADMIN_USERNAME")
	if adminUsername == "" {
		return nil, fmt.Errorf("ADMIN_USERNAME is not set")
	}

	adminPasswordHash := os.Getenv("ADMIN_PASSWORD_HASH")
	if adminPasswordHash == "" {
		return nil, fmt.Errorf("ADMIN_PASSWORD_HASH is not set")
	}

	jwtSecretRaw := os.Getenv("JWT_SECRET")
	if jwtSecretRaw == "" {
		return nil, fmt.Errorf("JWT_SECRET is not set")
	}

	ttlHoursRaw := os.Getenv("JWT_TOKEN_TTL_HOURS")
	if ttlHoursRaw == "" {
		return nil, fmt.Errorf("JWT_TOKEN_TTL_HOURS is not set")
	}
	ttlHours, err := strconv.Atoi(ttlHoursRaw)
	if err != nil {
		return nil, fmt.Errorf("parse JWT_TOKEN_TTL_HOURS: %w", err)
	}

	altLanguageLabel := os.Getenv("ALT_LANGUAGE_LABEL")
	if altLanguageLabel == "" {
		altLanguageLabel = "ALT" // sensible fallback if not configured
	}

	translationsDir := os.Getenv("TRANSLATIONS_DIR")
	if translationsDir == "" {
		translationsDir = "./locales"
	}

	adminCredentialsPath := os.Getenv("ADMIN_CREDENTIALS_PATH")
	if adminCredentialsPath == "" {
		adminCredentialsPath = "./admin_credentials.json"
	}

	uploadsDir := os.Getenv("UPLOADS_DIR")
	if uploadsDir == "" {
		uploadsDir = "./uploads"
	}

	deploymentMode := os.Getenv("DEPLOYMENT_MODE")
	if deploymentMode == "" {
		deploymentMode = "lan" // Default
	}

	return &Config{
		DBString:          dbString,
		AdminUsername:     adminUsername,
		AdminPasswordHash: adminPasswordHash,
		JWTSecret:         []byte(jwtSecretRaw),
		JWTTokenTTL:       time.Duration(ttlHours) * time.Hour,
		AltLanguageLabel:  altLanguageLabel,
		TranslationsDir:   translationsDir,
		UploadsDir:        uploadsDir,
		DeploymentMode:    deploymentMode,
		PublicBaseURL:     os.Getenv("PUBLIC_BASE_URL"),
	}, nil
}

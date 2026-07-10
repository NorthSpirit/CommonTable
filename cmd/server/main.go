package main

import (
	"context"
	"log"

	"CommonTable/internal/config"
	"CommonTable/internal/repository"
	"CommonTable/internal/services/adminservices"
	"CommonTable/internal/services/authservices"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	pool, err := pgxpool.New(context.Background(), cfg.DBString)
	if err != nil {
		log.Fatalf("connect to db: %v", err)
	}
	defer pool.Close()

	queries := repository.New(pool)

	authService := authservices.NewAuthService(queries, cfg.JWTSecret, cfg.JWTTokenTTL)

	adminService := adminservices.NewAdminService(
		adminservices.AdminConfig{
			Username:     cfg.AdminUsername,
			PasswordHash: cfg.AdminPasswordHash,
		},
		cfg.JWTSecret,
		cfg.JWTTokenTTL,
	)

	// ... wire up remaining services + handlers + router from here
	_ = authService
	_ = adminService
}

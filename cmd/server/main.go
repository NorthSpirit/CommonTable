package main

import (
	"CommonTable/internal/templates"
	"context"
	"log"
	"net/http"

	"CommonTable/internal/config"
	"CommonTable/internal/container"
	"CommonTable/internal/httpserver"
	"CommonTable/internal/repository"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}
	templates.SetAltLanguageLabel(cfg.AltLanguageLabel)
	pool, err := pgxpool.New(context.Background(), cfg.DBString)
	if err != nil {
		log.Fatalf("connect to db: %v", err)
	}
	defer pool.Close()

	queries := repository.New(pool)

	services := container.NewServices(queries, container.ServicesConfig{
		AdminUsername:     cfg.AdminUsername,
		AdminPasswordHash: cfg.AdminPasswordHash,
		JWTSecret:         cfg.JWTSecret,
		JWTTokenTTL:       cfg.JWTTokenTTL,
	})

	customerHandlers := container.NewCustomerHandlers(services)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	httpserver.RegisterCustomerRoutes(r, customerHandlers)
	r.Handle("/assets/*", http.StripPrefix("/assets", http.FileServer(http.Dir("./assets"))))

	log.Println("listening on :41241")
	if err := http.ListenAndServe(":41241", r); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

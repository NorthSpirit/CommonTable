package main

import (
	"CommonTable/internal/handlers/hub"
	"CommonTable/internal/templates/shared"
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
	shared.SetAltLanguageLabel(cfg.AltLanguageLabel)

	if err := shared.LoadTranslations(cfg.TranslationsDir); err != nil {
		log.Fatalf("load translations: %v", err)
	}

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
		UploadsDir:        cfg.UploadsDir,
	})

	customerHandlers := container.NewCustomerHandlers(services)
	adminHandlers := container.NewAdminHandlers(services)
	chefHandlers := container.NewChefHandlers(services)
	kitchenHandlers := container.NewKitchenHandlers(services)
	hubHandler := hub.NewHubHandler()

	waiterAPIHandlers := container.NewWaiterAPIHandlers(services)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	httpserver.RegisterCustomerRoutes(r, customerHandlers)
	httpserver.RegisterAdminRoutes(r, adminHandlers, services.Admin)
	httpserver.RegisterChefRoutes(r, chefHandlers, services.Auth)
	httpserver.RegisterKitchenRoutes(r, kitchenHandlers, services.Auth)

	httpserver.RegisterWaiterAPIRoutes(r, waiterAPIHandlers, services.Auth)

	r.Handle("/assets/*", http.StripPrefix("/assets", http.FileServer(http.Dir("./assets"))))
	r.Handle("/uploads/*", http.StripPrefix("/uploads", http.FileServer(http.Dir(cfg.UploadsDir))))
	r.Get("/hub", hubHandler.Page)

	log.Println("listening on :41241")
	if err := http.ListenAndServe(":41241", r); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

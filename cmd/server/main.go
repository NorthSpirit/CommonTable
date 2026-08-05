package main

import (
	"context"
	"log"
	"net/http"

	"CommonTable/internal/config"
	"CommonTable/internal/container"
	"CommonTable/internal/discovery"
	"CommonTable/internal/handlers/hub"
	"CommonTable/internal/httpserver"
	"CommonTable/internal/logging"
	"CommonTable/internal/repository"
	"CommonTable/internal/templates/shared"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

const httpPort = ":41241"

func main() {
	logWriter, logFile, err := logging.Setup("./logs")
	if err != nil {
		log.Fatalf("failed to set up logging: %v", err)
	}
	defer logFile.Close()

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
	adminHandlers := container.NewAdminHandlers(services, pool)
	chefHandlers := container.NewChefHandlers(services)
	kitchenHandlers := container.NewKitchenHandlers(services)
	waiterAPIHandlers := container.NewWaiterAPIHandlers(services, pool)
	hubHandler := hub.NewHubHandler()

	r := chi.NewRouter()
	r.Use(middleware.RequestLogger(&middleware.DefaultLogFormatter{
		Logger:  log.New(logWriter, "", log.LstdFlags),
		NoColor: true,
	}))

	httpserver.RegisterCustomerRoutes(r, customerHandlers)
	httpserver.RegisterAdminRoutes(r, adminHandlers, services.Admin)
	httpserver.RegisterChefRoutes(r, chefHandlers, services.Auth)
	httpserver.RegisterKitchenRoutes(r, kitchenHandlers, services.Auth)
	httpserver.RegisterWaiterAPIRoutes(r, waiterAPIHandlers, services.Auth)

	r.Handle("/assets/*", http.StripPrefix("/assets", http.FileServer(http.Dir("./assets"))))
	r.Handle("/uploads/*", http.StripPrefix("/uploads", http.FileServer(http.Dir(cfg.UploadsDir))))
	r.Get("/hub", hubHandler.Page)
	r.Get("/download/waiter-app.apk", hub.ServeAPK("./assets/downloads/waiter_app.apk"))
	r.Get("/api/discover", hub.Discover)

	// UDP LAN discovery — lets clients find this server without a
	// hardcoded IP. Irrelevant/unused once deployed behind a real
	// hostname, but harmless to leave running either way.
	discovery.StartResponder(41241)

	log.Printf("listening on %s", httpPort)
	if err := http.ListenAndServe(httpPort, r); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

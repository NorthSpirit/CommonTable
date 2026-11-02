package main

import (
	"context"
	"log"
	"net/http"

	"CommonTable/internal/config"
	"CommonTable/internal/container"
	"CommonTable/internal/discovery"
	"CommonTable/internal/handlers/hub"
	"CommonTable/internal/handlers/terminal"
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
	// Sets logging up by creating the file, with the date of the server start as a name (and also printing to the screen)
	logWriter, logFile, err := logging.Setup("./logs")
	if err != nil {
		log.Fatalf("failed to set up logging: %v", err)
	}
	defer logFile.Close()

	// Attempts to read the config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	// Sets the label of the alt language (shown on the language switch button, like "LT"), from the config
	shared.SetAltLanguageLabel(cfg.AltLanguageLabel)

	// Reads all the translation files (every .json in the translations folder, file name = language, en.json is required as the fallback)
	if err := shared.LoadTranslations(cfg.TranslationsDir); err != nil {
		log.Fatalf("load translations: %v", err)
	}

	// Opens a connection POOL to Postgres
	// context.Background() = a plain, never-cancelled context.
	pool, err := pgxpool.New(context.Background(), cfg.DBString)
	if err != nil {
		log.Fatalf("connect to db: %v", err)
	}
	defer pool.Close() // close all DB connections on shutdown

	// Queries in this project is SQLC auto-generated code
	// Each method uses one SQL query
	queries := repository.New(pool)

	// Business logic stuff
	// Builds all the business-logic services in one place (using dependency injection containers)
	// Services use the DB queries and get the setting they need: Admin login, JWT signing, upload folder.
	services := container.NewServices(queries, container.ServicesConfig{
		AdminUsername:     cfg.AdminUsername,
		AdminPasswordHash: cfg.AdminPasswordHash,
		JWTSecret:         cfg.JWTSecret,   // key used to sign login tokens
		JWTTokenTTL:       cfg.JWTTokenTTL, // how long a token stays valid
		UploadsDir:        cfg.UploadsDir,  // where uploaded files are stored
	})

	// Handler stuff (HTTP layers + waiter's app)
	// One set of handlers per user role / part of the app.
	// Each handler receives the services it needs.
	customerHandlers := container.NewCustomerHandlers(services)
	adminHandlers := container.NewAdminHandlers(services, pool)
	chefHandlers := container.NewChefHandlers(services)
	kitchenHandlers := container.NewKitchenHandlers(services)
	waiterAPIHandlers := container.NewWaiterAPIHandlers(services, pool)
	managerHandlers := container.NewManagerHandlers(services, pool)
	staffHandlers := container.NewStaffHandlers(services)

	// Hub is a "landing" page (links, APK download, API docs). It doesn't need services.
	hubHandler := hub.NewHubHandler()

	// Since terminal is a bit extra functionality, not really meant to be on the end product
	// It's handler is built here, to be, probably eventually deleted later on.
	terminalHandlers := httpserver.TerminalHandlers{Handler: terminal.NewHandler(services.Auth, services.FoodOrder, services.ItemOrder, services.Payment)}

	// Router and middleware stuff

	// chi is used for the router - it matches URL + method to a handler function.
	r := chi.NewRouter()

	// Middleware runs on every request, logging each request.
	r.Use(middleware.RequestLogger(&middleware.DefaultLogFormatter{
		Logger:  log.New(logWriter, "", log.LstdFlags),
		NoColor: true,
	}))

	// Catches panics inside handlers, so the user gets a proper 500 response and the error is logged, instead of the connection just being dropped
	r.Use(middleware.Recoverer)

	// Each Register function attaches that role's URLs to the router.
	// Most of them use service.Auth (or admin) to authenticate, before letting anybody through
	// Obviously, customers doesn't need authentication, cause rando customers doesn't need to log in.
	httpserver.RegisterCustomerRoutes(r, customerHandlers)
	httpserver.RegisterAdminRoutes(r, adminHandlers, services.Admin)
	httpserver.RegisterChefRoutes(r, chefHandlers, services.Auth)
	httpserver.RegisterKitchenRoutes(r, kitchenHandlers, services.Auth)
	httpserver.RegisterWaiterAPIRoutes(r, waiterAPIHandlers, services.Auth)
	httpserver.RegisterManagerRoutes(r, managerHandlers, services.Auth)
	httpserver.RegisterStaffRoutes(r, staffHandlers, services.Auth)
	httpserver.RegisterTerminalRoutes(r, terminalHandlers)

	// Static files. StripPrefix removes "/assets" from the URL before looking in the folder
	r.Handle("/assets/*", http.StripPrefix("/assets", http.FileServer(http.Dir("./assets"))))

	//Basically, same as with assets, removes a part from the URL.
	r.Handle("/uploads/*", http.StripPrefix("/uploads", http.FileServer(http.Dir(cfg.UploadsDir))))

	// Hub and it's related routes
	r.Get("/hub", hubHandler.Page)
	r.Get("/download/waiter-app.apk", hub.ServeAPK("./assets/downloads/waiter_app.apk"))
	r.Get("/api/discover", hub.Discover)
	r.Get("/api-docs", hubHandler.ApiDocs)
	r.Get("/api-docs/export", hubHandler.ExportDocs)

	// UDP LAN discovery — lets clients find this server without a hardcoded IP.
	// Only runs in a LAN mode
	// Hosted deployment (the one used for demo) has no use for it
	// and there's no local subnet to broadcast on anyhow.
	if cfg.DeploymentMode == "lan" {
		// Listens for UDP "where are you?" broadcasts and replies with the service name and HTTP port, runs in the background
		discovery.StartResponder(41241)
		log.Printf("listening on %s (LAN mode)", httpPort)
		discovery.PrintLANAddresses()
	} else {
		log.Printf("listening on %s (hosted mode, public URL: %s)", httpPort, cfg.PublicBaseURL)
	}

	// Finally starts the server
	if err := http.ListenAndServe(httpPort, r); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

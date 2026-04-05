// Migrated from: CoffeeDiaryApplication.java
package main

import (
	"context"
	"database/sql"
	"encoding/gob"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"

	"github.com/oglimmer/coffee-diary-backend/internal/config"
	"github.com/oglimmer/coffee-diary-backend/internal/handler"
	"github.com/oglimmer/coffee-diary-backend/internal/repository"
	"github.com/oglimmer/coffee-diary-backend/internal/service"
)

func main() {
	// Register int64 for gob encoding (gorilla/sessions stores values via gob)
	gob.Register(int64(0))

	cfg := config.Load()

	// Database
	db, err := sql.Open("mysql", cfg.DSN())
	if err != nil {
		slog.Error("failed to open database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		slog.Error("failed to ping database", "error", err)
		os.Exit(1)
	}

	// Run migrations
	if err := runMigrations(db); err != nil && err != migrate.ErrNoChange {
		slog.Error("failed to run migrations", "error", err)
		os.Exit(1)
	}

	// OIDC setup
	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, cfg.OIDCIssuerURL)
	if err != nil {
		slog.Error("failed to create OIDC provider", "error", err)
		os.Exit(1)
	}

	oauth2Config := &oauth2.Config{
		ClientID:     cfg.OIDCClientID,
		ClientSecret: cfg.OIDCClientSecret,
		RedirectURL:  cfg.OIDCRedirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	verifier := provider.Verifier(&oidc.Config{ClientID: cfg.OIDCClientID})

	// Apple Sign in with Apple OIDC provider
	appleProvider, err := oidc.NewProvider(ctx, "https://appleid.apple.com")
	if err != nil {
		slog.Error("failed to create Apple OIDC provider", "error", err)
		os.Exit(1)
	}
	appleVerifier := appleProvider.Verifier(&oidc.Config{ClientID: cfg.AppleClientID})

	// Discover end_session_endpoint from provider
	var providerClaims struct {
		EndSessionEndpoint string `json:"end_session_endpoint"`
	}
	if err := provider.Claims(&providerClaims); err != nil {
		slog.Warn("could not read end_session_endpoint from OIDC provider", "error", err)
	}

	// Session store
	store := sessions.NewCookieStore([]byte(cfg.SessionSecret))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,
		Secure:   cfg.CookieSecure,
		SameSite: http.SameSiteLaxMode,
	}

	// Repositories
	userRepo := repository.NewUserRepository(db)
	coffeeRepo := repository.NewCoffeeRepository(db)
	sieveRepo := repository.NewSieveRepository(db)
	entryRepo := repository.NewDiaryEntryRepository(db)

	// Services
	appleTokenSvc := service.NewAppleTokenService(cfg.AppleClientID, cfg.AppleTeamID, cfg.AppleKeyID, cfg.ApplePrivateKey)
	authSvc := service.NewAuthService(userRepo, appleTokenSvc)
	coffeeSvc := service.NewCoffeeService(coffeeRepo)
	sieveSvc := service.NewSieveService(sieveRepo)
	entrySvc := service.NewDiaryEntryService(entryRepo, coffeeRepo, sieveRepo)

	// Handlers
	authH := handler.NewAuthHandler(authSvc, appleTokenSvc, oauth2Config, verifier, appleVerifier, store, cfg.FrontendURL, providerClaims.EndSessionEndpoint)
	coffeeH := handler.NewCoffeeHandler(coffeeSvc)
	sieveH := handler.NewSieveHandler(sieveSvc)
	entryH := handler.NewDiaryEntryHandler(entrySvc)
	actuatorH := handler.NewActuatorHandler(cfg)

	// Router
	r := chi.NewRouter()
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(handler.SecurityHeaders)

	// Auth endpoints
	r.Route("/api/auth", func(r chi.Router) {
		// Public
		r.Get("/login", authH.Login)
		r.Get("/callback", authH.Callback)
		r.Get("/logout", authH.Logout)
		r.Post("/apple-callback", authH.AppleCallback)
		r.Get("/me", authH.Me)
		// Protected: account deletion
		r.With(handler.SessionAuth(store)).Delete("/me", authH.DeleteAccount)
	})

	// Protected API endpoints
	r.Group(func(r chi.Router) {
		r.Use(handler.SessionAuth(store))

		r.Route("/api/coffees", func(r chi.Router) {
			r.Get("/", coffeeH.FindAll)
			r.Post("/", coffeeH.Create)
			r.Delete("/{id}", coffeeH.Delete)
		})

		r.Route("/api/sieves", func(r chi.Router) {
			r.Get("/", sieveH.FindAll)
			r.Post("/", sieveH.Create)
			r.Delete("/{id}", sieveH.Delete)
		})

		r.Route("/api/diary-entries", func(r chi.Router) {
			r.Get("/", entryH.FindAll)
			r.Get("/{id}", entryH.FindByID)
			r.Post("/", entryH.Create)
			r.Put("/{id}", entryH.Update)
			r.Delete("/{id}", entryH.Delete)
		})
	})

	// Actuator endpoints
	r.Route("/actuator", func(r chi.Router) {
		r.Get("/health", actuatorH.Health)
		r.Get("/info", actuatorH.Info)
		r.Handle("/prometheus", actuatorH.Prometheus())
		r.Get("/metrics", actuatorH.Metrics)
	})

	// Server with graceful shutdown
	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: r,
	}

	go func() {
		slog.Info("starting server", "port", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("shutting down server")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("server shutdown error", "error", err)
	}
}

func runMigrations(db *sql.DB) error {
	// If migrating from Flyway: seed golang-migrate's tracking table so it
	// knows the baseline schema (migration 1) is already applied.
	if flywayExists(db) && !golangMigrateExists(db) {
		slog.Info("detected Flyway history — seeding golang-migrate baseline")
		if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (version bigint NOT NULL PRIMARY KEY, dirty boolean NOT NULL)`); err != nil {
			return fmt.Errorf("create schema_migrations: %w", err)
		}
		if _, err := db.Exec(`INSERT INTO schema_migrations (version, dirty) VALUES (1, false)`); err != nil {
			return fmt.Errorf("seed schema_migrations: %w", err)
		}
	}

	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return fmt.Errorf("migration driver: %w", err)
	}
	m, err := migrate.NewWithDatabaseInstance("file://migrations", "mysql", driver)
	if err != nil {
		return fmt.Errorf("migration init: %w", err)
	}
	return m.Up()
}

func flywayExists(db *sql.DB) bool {
	var n int
	err := db.QueryRow("SELECT 1 FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = 'flyway_schema_history' LIMIT 1").Scan(&n)
	return err == nil
}

func golangMigrateExists(db *sql.DB) bool {
	var n int
	err := db.QueryRow("SELECT 1 FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = 'schema_migrations' LIMIT 1").Scan(&n)
	return err == nil
}

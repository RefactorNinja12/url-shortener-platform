package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/rini/url-shortener/go-api/internal/handler"
	"github.com/rini/url-shortener/go-api/internal/queue"
	"github.com/rini/url-shortener/go-api/internal/store"
	"github.com/rini/url-shortener/go-api/migrations"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	dbURL := envOrDefault("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/urlshortener?sslmode=disable")

	if err := runMigrations(dbURL); err != nil {
		logger.Error("migration failed", "error", err)
		os.Exit(1)
	}

	db, err := store.New(context.Background(), dbURL)
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	rabbitURL := envOrDefault("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")
	publisher := queue.Connect(rabbitURL, logger)
	defer publisher.Close()

	baseURL := envOrDefault("BASE_URL", "http://localhost:8080")
	h := handler.New(db, logger, baseURL, publisher)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("POST /shorten", h.ShortenURL)
	mux.HandleFunc("GET /{code}", h.Redirect)

	addr := envOrDefault("ADDR", ":8080")
	srv := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		logger.Info("server starting", "addr", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("graceful shutdown failed", "error", err)
	}
	logger.Info("server stopped")
}

func runMigrations(dbURL string) error {
	// golang-migrate:s pgx/v5-driver kräver pgx5:// som schema istället för postgres://
	migrateURL := strings.Replace(dbURL, "postgres://", "pgx5://", 1)

	src, err := iofs.New(migrations.FS, ".")
	if err != nil {
		return err
	}

	m, err := migrate.NewWithSourceInstance("iofs", src, migrateURL)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

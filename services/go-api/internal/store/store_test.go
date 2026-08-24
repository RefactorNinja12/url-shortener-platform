package store_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"

	"github.com/rini/url-shortener/go-api/internal/store"
	"github.com/rini/url-shortener/go-api/migrations"
)

// setupTestStore startar en riktig Postgres i en container, kör migrationerna
// mot den och returnerar en anslutet Store. Kräver tillgång till Docker.
func setupTestStore(t *testing.T) *store.Store {
	t.Helper()
	ctx := context.Background()

	pgContainer, err := tcpostgres.Run(ctx, "postgres:16-alpine",
		tcpostgres.WithDatabase("urlshortener"),
		tcpostgres.WithUsername("postgres"),
		tcpostgres.WithPassword("postgres"),
		tcpostgres.BasicWaitStrategies(),
	)
	if err != nil {
		t.Fatalf("failed to start postgres container: %v", err)
	}
	t.Cleanup(func() {
		if err := pgContainer.Terminate(context.Background()); err != nil {
			t.Logf("failed to terminate postgres container: %v", err)
		}
	})

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("failed to get connection string: %v", err)
	}

	if err := waitForReady(ctx, connStr); err != nil {
		t.Fatalf("postgres never became ready: %v", err)
	}

	if err := runMigrations(connStr); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	s, err := store.New(ctx, connStr)
	if err != nil {
		t.Fatalf("failed to connect store: %v", err)
	}
	t.Cleanup(s.Close)

	return s
}

func waitForReady(ctx context.Context, connStr string) error {
	deadline := time.Now().Add(30 * time.Second)
	var lastErr error
	for time.Now().Before(deadline) {
		s, err := store.New(ctx, connStr)
		if err == nil {
			s.Close()
			return nil
		}
		lastErr = err
		time.Sleep(500 * time.Millisecond)
	}
	return lastErr
}

func runMigrations(connStr string) error {
	migrateURL := strings.Replace(connStr, "postgres://", "pgx5://", 1)

	src, err := iofs.New(migrations.FS, ".")
	if err != nil {
		return err
	}

	m, err := migrate.NewWithSourceInstance("iofs", src, migrateURL)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}

func TestStore_CreateAndGetURL(t *testing.T) {
	s := setupTestStore(t)
	ctx := context.Background()

	ownerID := int64(7)
	created, err := s.CreateURL(ctx, "abc1234", "https://example.com", &ownerID)
	if err != nil {
		t.Fatalf("CreateURL failed: %v", err)
	}
	if created.Code != "abc1234" {
		t.Errorf("expected code %q, got %q", "abc1234", created.Code)
	}
	if created.OriginalURL != "https://example.com" {
		t.Errorf("expected original URL %q, got %q", "https://example.com", created.OriginalURL)
	}
	if created.ID == 0 {
		t.Error("expected a non-zero ID")
	}
	if created.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
	if created.OwnerID == nil || *created.OwnerID != ownerID {
		t.Errorf("expected OwnerID %d, got %v", ownerID, created.OwnerID)
	}

	fetched, err := s.GetURL(ctx, "abc1234")
	if err != nil {
		t.Fatalf("GetURL failed: %v", err)
	}
	if fetched.OriginalURL != created.OriginalURL {
		t.Errorf("fetched URL %q does not match created URL %q", fetched.OriginalURL, created.OriginalURL)
	}
}

func TestStore_GetURL_NotFound(t *testing.T) {
	s := setupTestStore(t)
	ctx := context.Background()

	_, err := s.GetURL(ctx, "doesnotexist")
	if !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestStore_CreateURL_DuplicateCode(t *testing.T) {
	s := setupTestStore(t)
	ctx := context.Background()

	if _, err := s.CreateURL(ctx, "dupe123", "https://example.com/one", nil); err != nil {
		t.Fatalf("first CreateURL failed: %v", err)
	}

	_, err := s.CreateURL(ctx, "dupe123", "https://example.com/two", nil)
	if err == nil {
		t.Fatal("expected an error when inserting a duplicate code, got nil")
	}
}

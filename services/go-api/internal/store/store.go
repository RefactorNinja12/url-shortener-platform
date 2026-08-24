package store

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("not found")

type Store struct {
	pool *pgxpool.Pool
}

type URL struct {
	ID          int64
	Code        string
	OriginalURL string
	CreatedAt   time.Time
	OwnerID     *int64
}

func New(ctx context.Context, connString string) (*Store, error) {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	return &Store{pool: pool}, nil
}

func (s *Store) Close() {
	s.pool.Close()
}

// CreateURL sparar en ny URL-mappning och returnerar den sparade posten.
func (s *Store) CreateURL(ctx context.Context, code, originalURL string) (*URL, error) {
	const q = `INSERT INTO urls (code, original_url) VALUES ($1, $2)
	           RETURNING id, code, original_url, created_at, owner_id`

	var u URL
	row := s.pool.QueryRow(ctx, q, code, originalURL)
	if err := row.Scan(&u.ID, &u.Code, &u.OriginalURL, &u.CreatedAt, &u.OwnerID); err != nil {
		return nil, err
	}
	return &u, nil
}

// GetURL hämtar en URL-mappning baserat på kort-kod.
// Returnerar ErrNotFound om koden inte finns.
func (s *Store) GetURL(ctx context.Context, code string) (*URL, error) {
	const q = `SELECT id, code, original_url, created_at, owner_id
	           FROM urls WHERE code = $1`

	var u URL
	row := s.pool.QueryRow(ctx, q, code)
	if err := row.Scan(&u.ID, &u.Code, &u.OriginalURL, &u.CreatedAt, &u.OwnerID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &u, nil
}

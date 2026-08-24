package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/rini/url-shortener/go-api/internal/shortener"
	"github.com/rini/url-shortener/go-api/internal/store"
)

const maxCodeGenerationAttempts = 5

const pgUniqueViolation = "23505"

type Handler struct {
	store   *store.Store
	logger  *slog.Logger
	baseURL string
}

func New(s *store.Store, logger *slog.Logger, baseURL string) *Handler {
	return &Handler{store: s, logger: logger, baseURL: baseURL}
}

type shortenRequest struct {
	URL string `json:"url"`
}

type shortenResponse struct {
	ShortURL string `json:"shortUrl"`
}

func (h *Handler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	var req shortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if !isValidURL(req.URL) {
		http.Error(w, "url must be a valid absolute http(s) URL", http.StatusBadRequest)
		return
	}

	u, err := h.createWithUniqueCode(r.Context(), req.URL)
	if err != nil {
		h.logger.Error("failed to create short url", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, shortenResponse{ShortURL: h.baseURL + "/" + u.Code})
}

// createWithUniqueCode genererar en kod och försöker spara den. Vid en
// kollision (UNIQUE-constraint) genereras en ny kod och vi försöker igen,
// upp till maxCodeGenerationAttempts gånger.
func (h *Handler) createWithUniqueCode(ctx context.Context, originalURL string) (*store.URL, error) {
	var lastErr error
	for attempt := 0; attempt < maxCodeGenerationAttempts; attempt++ {
		code, err := shortener.GenerateCode()
		if err != nil {
			return nil, err
		}

		u, err := h.store.CreateURL(ctx, code, originalURL)
		if err == nil {
			return u, nil
		}
		if !isUniqueViolation(err) {
			return nil, err
		}
		lastErr = err
	}
	return nil, lastErr
}

func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code") // Go 1.22+: hämtar {code} ur URL-mönstret

	u, err := h.store.GetURL(r.Context(), code)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			http.NotFound(w, r)
			return
		}
		h.logger.Error("failed to look up code", "code", code, "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// TODO (Fas 5): publicera klick-event till kön här.

	http.Redirect(w, r, u.OriginalURL, http.StatusFound)
}

func isValidURL(raw string) bool {
	if raw == "" {
		return false
	}
	parsed, err := url.ParseRequestURI(raw)
	if err != nil {
		return false
	}
	if parsed.Host == "" {
		return false
	}
	return parsed.Scheme == "http" || parsed.Scheme == "https"
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == pgUniqueViolation
	}
	return false
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const testSecret = "test-secret"

func signToken(t *testing.T, secret string, claims jwt.MapClaims) string {
	t.Helper()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("failed to sign test token: %v", err)
	}
	return signed
}

func requestWithAuth(header string) *http.Request {
	r := httptest.NewRequest(http.MethodPost, "/shorten", nil)
	if header != "" {
		r.Header.Set("Authorization", header)
	}
	return r
}

func TestUserIDFromRequest_ValidToken(t *testing.T) {
	v := NewVerifier(testSecret)
	token := signToken(t, testSecret, jwt.MapClaims{
		"sub":   "42",
		"email": "test@example.com",
		"exp":   time.Now().Add(time.Hour).Unix(),
	})

	userID, err := v.UserIDFromRequest(requestWithAuth("Bearer " + token))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if userID != 42 {
		t.Errorf("expected userID 42, got %d", userID)
	}
}

func TestUserIDFromRequest_MissingHeader(t *testing.T) {
	v := NewVerifier(testSecret)
	_, err := v.UserIDFromRequest(requestWithAuth(""))
	if err != ErrMissingToken {
		t.Errorf("expected ErrMissingToken, got %v", err)
	}
}

func TestUserIDFromRequest_WrongSecret(t *testing.T) {
	v := NewVerifier(testSecret)
	token := signToken(t, "wrong-secret", jwt.MapClaims{
		"sub": "42",
		"exp": time.Now().Add(time.Hour).Unix(),
	})

	_, err := v.UserIDFromRequest(requestWithAuth("Bearer " + token))
	if err != ErrInvalidToken {
		t.Errorf("expected ErrInvalidToken, got %v", err)
	}
}

func TestUserIDFromRequest_ExpiredToken(t *testing.T) {
	v := NewVerifier(testSecret)
	token := signToken(t, testSecret, jwt.MapClaims{
		"sub": "42",
		"exp": time.Now().Add(-time.Hour).Unix(),
	})

	_, err := v.UserIDFromRequest(requestWithAuth("Bearer " + token))
	if err != ErrInvalidToken {
		t.Errorf("expected ErrInvalidToken, got %v", err)
	}
}

func TestUserIDFromRequest_MalformedBearer(t *testing.T) {
	v := NewVerifier(testSecret)
	_, err := v.UserIDFromRequest(requestWithAuth("Basic sometoken"))
	if err != ErrMissingToken {
		t.Errorf("expected ErrMissingToken, got %v", err)
	}
}

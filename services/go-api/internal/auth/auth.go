package auth

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrMissingToken = errors.New("missing bearer token")
	ErrInvalidToken = errors.New("invalid token")
)

// Verifier validerar JWT:er utfärdade av auth-admin, som signerar med samma
// delade hemlighet (HMAC).
type Verifier struct {
	secret []byte
}

func NewVerifier(secret string) *Verifier {
	return &Verifier{secret: []byte(secret)}
}

// UserIDFromRequest läser Authorization-headern, validerar Bearer-token:et
// och returnerar användar-ID:t ur "sub"-claimet.
func (v *Verifier) UserIDFromRequest(r *http.Request) (int64, error) {
	header := r.Header.Get("Authorization")
	const prefix = "Bearer "
	if !strings.HasPrefix(header, prefix) {
		return 0, ErrMissingToken
	}
	tokenString := strings.TrimPrefix(header, prefix)

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return v.secret, nil
	})
	if err != nil || !token.Valid {
		return 0, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, ErrInvalidToken
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return 0, ErrInvalidToken
	}

	userID, err := strconv.ParseInt(sub, 10, 64)
	if err != nil {
		return 0, ErrInvalidToken
	}

	return userID, nil
}

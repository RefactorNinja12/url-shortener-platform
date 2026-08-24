package shortener

import (
	"crypto/rand"
	"math/big"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// GenerateCode returnerar en kryptografiskt slumpmässig base62-sträng.
// Längden 7 ger 62^7 ≈ 3,5 biljoner möjliga koder — gott om utrymme.
//
// Alternativ att fundera på (se "Öppna beslut" i projektplanen):
//   - base62-encoding av auto-increment ID undviker kollisioner helt
//     men kräver att du hämtar ID:t från DB innan du genererar koden.
//   - Slumpmässig kod (som här) är enklare men kräver att du hanterar
//     UNIQUE-constraint-fel vid INSERT och försöker igen.
func GenerateCode() (string, error) {
	code := make([]byte, 7)
	for i := range code {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		if err != nil {
			return "", err
		}
		code[i] = alphabet[n.Int64()]
	}
	return string(code), nil
}

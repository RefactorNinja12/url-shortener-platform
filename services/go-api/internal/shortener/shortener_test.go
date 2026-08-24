package shortener

import (
	"strings"
	"testing"
)

func TestGenerateCode_LengthAndAlphabet(t *testing.T) {
	code, err := GenerateCode()
	if err != nil {
		t.Fatalf("GenerateCode returned error: %v", err)
	}
	if len(code) != 7 {
		t.Fatalf("expected code length 7, got %d (%q)", len(code), code)
	}
	for _, c := range code {
		if !strings.ContainsRune(alphabet, c) {
			t.Fatalf("code %q contains character %q not in alphabet", code, c)
		}
	}
}

func TestGenerateCode_LooksRandom(t *testing.T) {
	seen := make(map[string]bool)
	const n = 1000
	for i := 0; i < n; i++ {
		code, err := GenerateCode()
		if err != nil {
			t.Fatalf("GenerateCode returned error: %v", err)
		}
		if seen[code] {
			t.Fatalf("got duplicate code %q within %d generations", code, n)
		}
		seen[code] = true
	}
}

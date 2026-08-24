package handler

import "testing"

func TestIsValidURL(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want bool
	}{
		{"empty", "", false},
		{"valid https", "https://example.com/path?x=1", true},
		{"valid http", "http://example.com", true},
		{"no scheme", "example.com", false},
		{"unsupported scheme", "ftp://example.com", false},
		{"not a url", "not a url", false},
		{"scheme only, no host", "https://", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := isValidURL(tc.in)
			if got != tc.want {
				t.Errorf("isValidURL(%q) = %v, want %v", tc.in, got, tc.want)
			}
		})
	}
}

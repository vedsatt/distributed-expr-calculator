package jwt

import (
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestGenerate(t *testing.T) {
	token, err := Generate(123)
	if err != nil {
		t.Fatalf("Generte error: %v", err)
	}

	if token == "" {
		t.Fatalf("Generation token fail: empty token")
	}

	token2, _ := Generate(321)
	if token == token2 {
		t.Fatalf("Generate returns same tokens for different urernames")
	}

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		t.Errorf("Invalid token format, got %d parts", len(parts))
	}
}

func TestVerify(t *testing.T) {
	testcases := []struct {
		name     string
		id       int
		token    string
		expected bool
	}{
		{
			name:     "valid",
			id:       123,
			token:    GeneratetestTokens(t, 123, time.Now().Add(time.Minute*5)),
			expected: true,
		},
		{
			name:     "invalid",
			token:    "better.call.saul",
			expected: false,
		},
		{
			name:     "expired",
			id:       123,
			token:    GeneratetestTokens(t, 123, time.Now().Add(-time.Minute*5)),
			expected: false,
		},
	}

	for _, tc := range testcases {
		got, id := Verify(tc.token)
		if tc.expected != got {
			t.Fatalf("test: %v failure, expected: %v, but got: %v", tc.name, tc.expected, got)
		}
		if tc.id != id && got {
			t.Fatalf("test: %v failure id, expected: %v, but got: %v", tc.name, tc.id, id)
		}
	}
}

func GeneratetestTokens(t *testing.T, userID int, exp time.Time) string {
	t.Helper()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  userID,
		"exp": exp.Unix(),
		"iat": time.Now().Unix(),
	})

	tokenString, _ := token.SignedString([]byte(HmacSampleSecret))
	return tokenString
}

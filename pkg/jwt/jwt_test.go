package jwt

import (
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestGenerate(t *testing.T) {
	username1 := "username1"

	token, err := Generate(username1)
	if err != nil {
		t.Fatalf("Generte error: %v", err)
	}

	if token == "" {
		t.Fatalf("Generation token fail: empty token")
	}

	username2 := "username2"
	token2, _ := Generate(username2)
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
		token    string
		expected bool
	}{
		{
			name:     "valid",
			token:    GeneratetestTokens(t, "usname", time.Now().Add(time.Minute*5)),
			expected: true,
		},
		{
			name:     "invalid",
			token:    "better.call.saul",
			expected: false,
		},
		{
			name:     "expired",
			token:    GeneratetestTokens(t, "usname", time.Now().Add(-time.Minute*5)),
			expected: false,
		},
	}

	for _, tc := range testcases {
		got := Verify(tc.token)
		if tc.expected != got {
			t.Fatalf("test: %v failure, expected: %v, but got: %v", tc.name, tc.expected, got)
		}
	}
}

func GeneratetestTokens(t *testing.T, usname string, exp time.Time) string {
	t.Helper()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name": usname,
		"exp":  exp.Unix(),
		"iat":  time.Now().Unix(),
	})

	tokenString, _ := token.SignedString([]byte(HmacSampleSecret))
	return tokenString
}

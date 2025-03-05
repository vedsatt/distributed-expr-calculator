package ast

import (
	"testing"
)

func TestTokens(t *testing.T) {
	tests := []struct {
		input    string
		expected []*token
	}{
		{
			input: "2+3",
			expected: []*token{
				{t: operand, val: "2"},
				{t: operator, val: "+"},
				{t: operand, val: "3"},
			},
		},
		{
			input: "2*(3+4)",
			expected: []*token{
				{t: operand, val: "2"},
				{t: operator, val: "*"},
				{t: openBracket, val: "("},
				{t: operand, val: "3"},
				{t: operator, val: "+"},
				{t: operand, val: "4"},
				{t: closeBracket, val: ")"},
			},
		},
	}

	for _, tt := range tests {
		result := tokens(tt.input)
		if len(result) != len(tt.expected) {
			t.Errorf("tokens(%s) = %v, expected %v", tt.input, result, tt.expected)
		}
		for i := range result {
			if result[i].t != tt.expected[i].t || result[i].val != tt.expected[i].val {
				t.Errorf("tokens(%s) = %v, expected %v", tt.input, result, tt.expected)
			}
		}
	}
}

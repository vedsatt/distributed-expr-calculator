package ast

import (
	"testing"
)

func TestRpn(t *testing.T) {
	tests := []struct {
		tokens   []*token
		expected []*token
		err      error
	}{
		{
			tokens: []*token{
				{t: operand, val: "2"},
				{t: operand, val: "3"},
				{t: operator, val: "+"},
			},
			expected: []*token{
				{t: operand, val: "2"},
				{t: operand, val: "3"},
				{t: operator, val: "+"},
			},
			err: nil,
		},
	}

	for _, tt := range tests {
		result, err := rpn(tt.tokens)
		if err != tt.err {
			t.Errorf("rpn() error = %v, expected %v", err, tt.err)
		}
		if len(result) != len(tt.expected) {
			t.Errorf("rpn() = %v, expected %v", result, tt.expected)
		}
		for i := range result {
			if result[i].t != tt.expected[i].t || result[i].val != tt.expected[i].val {
				t.Errorf("rpn() = %v, expected %v", result, tt.expected)
			}
		}
	}
}

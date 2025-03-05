package ast

import (
	"testing"
)

func TestExpErr(t *testing.T) {
	tests := []struct {
		expression string
		err        error
	}{
		{"2+3", nil},
		{"2+", ErrInvalidExpression},
		{"+2", ErrOperatorFirst},
		{"2++3", ErrMergedOperators},
		{"2+()", ErrEmptyBrackets},
		{"2+)3", ErrNotOpenedBracket},
		{"2+3a", ErrWrongCharacter},
		{"2/0", ErrDivisionByZero},
		{"(", ErrInvalidExpression},
		{"", ErrNoOperators},
	}

	for _, tt := range tests {
		err := expErr(tt.expression)
		if err != tt.err {
			t.Errorf("expErr(%s) = %v, expected %v", tt.expression, err, tt.err)
		}
	}
}

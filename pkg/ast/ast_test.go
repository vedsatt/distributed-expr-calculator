package ast

import (
	"testing"

	"github.com/vedsatt/calc_prl/internal/models"
)

func TestPriority(t *testing.T) {
	tests := []struct {
		op       string
		expected int
		err      error
	}{
		{"/", 3, nil},
		{"*", 3, nil},
		{"+", 2, nil},
		{"-", 2, nil},
		{"(", 1, nil},
		{"?", 0, ErrUnknownOperator},
	}

	for _, tt := range tests {
		result, err := priority(tt.op)
		if result != tt.expected || err != tt.err {
			t.Errorf("priority(%s) = (%d, %v), expected (%d, %v)", tt.op, result, err, tt.expected, tt.err)
		}
	}
}

func compareAstNodes(a, b *models.AstNode) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.AstType == b.AstType &&
		a.Value == b.Value &&
		compareAstNodes(a.Left, b.Left) &&
		compareAstNodes(a.Right, b.Right)
}

func TestAst(t *testing.T) {
	tests := []struct {
		name     string
		tokens   []*token
		expected *models.AstNode
		err      error
	}{
		{
			name: "simple addition",
			tokens: []*token{
				{t: operand, val: "2"},
				{t: operand, val: "3"},
				{t: operator, val: "+"},
			},
			expected: &models.AstNode{
				AstType: "operation",
				Value:   "+",
				Left:    &models.AstNode{AstType: "number", Value: "2"},
				Right:   &models.AstNode{AstType: "number", Value: "3"},
			},
			err: nil,
		},
		{
			name: "invalid expression (not enough operands)",
			tokens: []*token{
				{t: operand, val: "2"},
				{t: operator, val: "+"},
			},
			expected: nil,
			err:      ErrInvalidExpression,
		},
		{
			name: "complex expression",
			tokens: []*token{
				{t: operand, val: "2"},
				{t: operand, val: "3"},
				{t: operator, val: "*"},
				{t: operand, val: "4"},
				{t: operator, val: "+"},
			},
			expected: &models.AstNode{
				AstType: "operation",
				Value:   "+",
				Left: &models.AstNode{
					AstType: "operation",
					Value:   "*",
					Left:    &models.AstNode{AstType: "number", Value: "2"},
					Right:   &models.AstNode{AstType: "number", Value: "3"},
				},
				Right: &models.AstNode{AstType: "number", Value: "4"},
			},
			err: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ast(tt.tokens)
			if err != tt.err {
				t.Errorf("ast() error = %v, expected %v", err, tt.err)
			}
			if !compareAstNodes(result, tt.expected) {
				t.Errorf("ast() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

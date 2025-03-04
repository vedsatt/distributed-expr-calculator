package ast

import (
	"time"
)

type AstNode struct {
	ID       int           `json:"id"`
	AstType  string        `json:"type"`
	Value    string        `json:"operation"`
	Left     *AstNode      `json:"arg1"`
	Right    *AstNode      `json:"arg2"`
	Counting bool          `json:"status"`
	OpTime   time.Duration `json:"operation_time"`
}

func priority(op string) (int, error) {
	switch {
	case op == "/" || op == "*":
		return 3, nil
	case op == "+" || op == "-":
		return 2, nil
	case op == "(":
		return 1, nil
	default:
		return 0, ErrUnknownOperator
	}
}

func createNode(id int, val string, left, right *AstNode) *AstNode {
	node := &AstNode{
		ID:      id,
		AstType: "operation",
		Value:   val,
		Left:    left,
		Right:   right,
	}

	switch val {
	case "*":
		node.OpTime = TIME_MULTIPLICATIONS_MS
	case "/":
		node.OpTime = TIME_DIVISIONS_MS
	case "+":
		node.OpTime = TIME_ADDITION_MS
	case "-":
		node.OpTime = TIME_SUBTRACTION_MS
	}

	return node
}

func ast(tokens []*token) (*AstNode, error) {
	var stack []*AstNode
	id := 0

	for _, tok := range tokens {
		switch tok.t {
		case operand:
			// создаем узел для числа
			node := &AstNode{
				ID:      id,
				AstType: "number",
				Value:   tok.val,
			}
			stack = append(stack, node)
			id++

		case operator:
			// один оператор - два операнда
			if len(stack) < 2 {
				return nil, ErrInvalidExpression
			}

			// извлекаем правый и левый операнды (порядок важен)
			right := stack[len(stack)-1]
			left := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			// создаем новый узел операции для оператора
			node := createNode(id, tok.val, left, right)
			stack = append(stack, node)
			id++

		default:
			return nil, ErrWrongCharacter
		}
	}

	if len(stack) != 1 {
		return nil, ErrInvalidExpression
	}

	return stack[0], nil
}

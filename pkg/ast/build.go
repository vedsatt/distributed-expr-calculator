package ast

// модуль, преобразующий строку-выражение в ast

import (
	"strings"
)

func Build(expression string) (*AstNode, error) {
	expression = strings.ReplaceAll(expression, " ", "") // избавляемся от пробелов

	err := expErr(expression)
	if err != nil {
		return nil, err
	}

	tokens := tokens(expression)

	rpn, err := rpn(tokens)
	if err != nil {
		return nil, err
	}

	astRoot, err := ast(rpn)
	if err != nil {
		return nil, err
	}

	return astRoot, nil
}

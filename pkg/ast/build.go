package ast

// модуль, преобразующий строку-выражение в ast

import (
	"strings"
	"sync"

	"github.com/vedsatt/calc_prl/internal/models"
)

var (
	mu sync.Mutex
)

func Build(expression string) (*models.AstNode, error) {
	mu.Lock()
	defer mu.Unlock()
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

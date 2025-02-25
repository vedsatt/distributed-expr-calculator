package ast

// модуль, преобразующий строку-выражение в ast дерево

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

type astNode struct {
	astType  string
	operator string
	left     *astNode
	right    *astNode
}

type stack []*token

func (s *stack) push(t *token) {
	*s = append(*s, t)
}

func (s *stack) pop() (*token, error) {
	if len(*s) == 0 {
		return nil, ErrEmptyStack
	}
	t := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return t, nil
}

func (s *stack) peek() *token {
	if len(*s) == 0 {
		return nil
	}
	return (*s)[len(*s)-1]
}

func (s *stack) len() int {
	return len(*s)
}

// структура для первоначального разбиения строки на токены
type token struct {
	t   string // тип токена
	val string // значение токена
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

// работает на регулярках, проверяет, является символ оператором
func typeCheck(symbol string) bool {
	pattern := "[+\\-*/]"
	r, err := regexp.Compile(pattern)
	if err != nil {
		log.Fatalf("./pkg/ast (func: typeCheck); regexp error occured: %s", err)
	}

	return r.MatchString(string(symbol))
}

func ast(tokens []*token) (*astNode, error) {
	var stack []*astNode

	for _, tok := range tokens {
		switch tok.t {
		case operand:
			// создаем узел для числа
			node := &astNode{
				astType:  "number",
				operator: tok.val,
			}
			stack = append(stack, node)

		case operator:
			// один оператор - два операнда
			if len(stack) < 2 {
				return nil, ErrInvalidExpression
			}

			// извлекаем правый и левый операнды (порядок важен)
			right := stack[len(stack)-1]
			left := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			// создаем новый узел операции
			node := &astNode{
				astType:  "operation",
				operator: tok.val,
				left:     left,
				right:    right,
			}
			stack = append(stack, node)

		default:
			return nil, ErrWrongCharacter
		}
	}

	if len(stack) != 1 {
		return nil, ErrInvalidExpression
	}

	return stack[0], nil
}

func rpn(tokens []*token) ([]*token, error) {
	var stack stack
	output := make([]*token, 0)

	for _, tok := range tokens {
		switch tok.t {
		case operand:
			output = append(output, tok)

		case operator:
			currPriority, err := priority(tok.val)
			if err != nil {
				return nil, err
			}

			// извлекаем операторы с большим или равным приоритетом
			for stack.len() > 0 {
				top := stack.peek()
				if top.t == openBracket {
					break // открывающая скобка прерывает извлечение
				}

				topPriority, err := priority(top.val)
				if err != nil {
					return nil, err
				}

				if topPriority >= currPriority {
					popped, _ := stack.pop()
					output = append(output, popped)
				} else {
					break
				}
			}
			stack.push(tok)

		case openBracket:
			stack.push(tok)

		case closeBracket:
			// извлекаем до открывающей скобки
			found := false
			for stack.len() > 0 {
				popped, err := stack.pop()
				if err != nil {
					return nil, ErrInvalidExpression
				}
				if popped.t == openBracket {
					found = true
					break
				}
				output = append(output, popped)
			}
			if !found {
				return nil, ErrNotOpenedBracket
			}

		default:
			return nil, ErrUnknownOperator
		}
	}

	// достаем оставшиеся операторы
	for stack.len() > 0 {
		popped, err := stack.pop()
		if err != nil {
			return nil, err
		}
		if popped.t == openBracket {
			return nil, ErrNotClosedBracket
		}
		output = append(output, popped)
	}

	return output, nil
}

func tokens(str string) []*token {
	tokens := make([]*token, 0)

	str = strings.ReplaceAll(str, " ", "") // избавляемся от пробелов
	i := 0
	for i < len(str) {
		switch {
		case typeCheck(string(str[i])): // если оператор
			tokens = append(tokens, &token{t: operator, val: string(str[i])})
			i++

		case str[i] >= 48 && str[i] <= 57: // если число
			tmp := ""
			for i < len(str) && ((str[i] >= 48 && str[i] <= 57) || str[i] == 44 || str[i] == 46) {
				tmp += string(str[i])
				i++
			}
			tokens = append(tokens, &token{t: operand, val: string(tmp)})

		case str[i] == 40 || str[i] == 41: // если скобка
			tp := openBracket
			if str[i] == 41 {
				tp = closeBracket
			}
			tokens = append(tokens, &token{t: tp, val: string(str[i])})
			i++

		default:
			i++
		}
	}

	return tokens
}

func pre_order(node *astNode) {
	if node != nil {
		fmt.Println(node)
		pre_order(node.left)
		pre_order(node.right)
	}
}

func Build(str string) error {
	tokens := tokens(str)
	rpn, err := rpn(tokens)
	if err != nil {
		return err
	}
	astRoot, err := ast(rpn)
	if err != nil {
		return err
	}

	pre_order(astRoot)

	return nil
}

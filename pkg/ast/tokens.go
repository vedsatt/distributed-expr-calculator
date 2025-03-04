package ast

import (
	"log"
	"regexp"
)

// структура для первоначального разбиения строки на токены
type token struct {
	t   string // тип токена
	val string // значение токена
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

func tokens(str string) []*token {
	tokens := make([]*token, 0)

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

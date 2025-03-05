package ast

// первоначальная проверка на ошибки
// понижает шанс пропустить ошибку в выражении
func expErr(expression string) error {
	len := len(expression)
	flag := false
	start := 0
	end := 0

	for i := 0; i < len; i++ {
		curr := expression[i]
		next := byte(0)
		if i < len-1 {
			next = expression[i+1]
		}

		if curr == '(' {
			start++
		}
		if curr == ')' {
			end++
		}
		if 48 <= curr && curr <= 57 && !flag {
			flag = true
		}

		switch {
		case i == 0 && (curr == ')' || curr == '*' || curr == '+' || curr == '-' || curr == '/'):
			return ErrOperatorFirst
		case i == len-1 && (curr == '*' || curr == '+' || curr == '-' || curr == '/'):
			return ErrOperatorLast
		case curr == '(' && next == ')':
			return ErrEmptyBrackets
		case curr == ')' && next == '(':
			return ErrMergedBrackets
		case (curr == '*' || curr == '+' || curr == '-' || curr == '/') && (next == '*' || next == '+' || next == '-' || next == '/'):
			return ErrMergedOperators
		case curr < '(' || curr > '9':
			return ErrWrongCharacter
		case len <= 2:
			return ErrInvalidExpression
		case curr == '/' && next == '0':
			return ErrDivisionByZero
		}
	}

	// базовая проверка на корректность скобок
	if start > end {
		return ErrNotClosedBracket
	} else if end > start {
		return ErrNotOpenedBracket
	}

	if !flag {
		return ErrNoOperators
	}
	return nil
}

package ast

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

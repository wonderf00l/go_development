/*
Пакет calc содержит:
1. stack.go/stack_test.go - реализаия стека и тесты
2. entities.go - основные сущности приложения
3. checker.go/checker_test.go - реализация методов, проверяющих входное выражение
4. printer.go - реализация вывода результата
5. calc.go/calc_test.go - основные функции для чтения выражения, парсинга и вычисления результата
*/
package calc

import (
	"bufio"
	"io"
	"strconv"
	"unicode"
)

var operPriority = map[string]int{
	"+": 0,
	"-": 0,
	"*": 1,
	"/": 1,
}

// чтение из источника --> проверка --> парсинг, вычисление --> вывод результата
func Calculator(in io.Reader, out io.Writer, config *Config) error {
	checker, printer := config.Checker, config.Printer
	scanner := bufio.NewScanner(in)

	for scanner.Scan() {
		expression := scanner.Text()
		if err := checker.Check(expression); err != nil {
			return err
		}
		result := ParseExpression(expression)
		if err := printer.Print(result, out); err != nil {
			return err
		}
	}
	return nil
}

// получение токена из выражения: числа, операторы
func GetNextToken(s string, currId *int) string {
	switch {
	case *currId > len(s):
		return ""
	case unicode.IsDigit(rune(s[*currId])) || s[*currId] == '.':
		var num string
		for ; *currId < len(s) && (unicode.IsDigit(rune(s[*currId])) || s[*currId] == '.'); *currId++ {
			num += string(s[*currId])
		}
		return num
	case s[*currId] == '(' && s[*currId+1] == '-':
		*currId += 2
		var negativeNum string
		for ; *currId < len(s) && unicode.IsDigit(rune(s[*currId])) || s[*currId] == '.'; *currId++ {
			negativeNum += string(s[*currId])
		}
		if *currId < len(s) && s[*currId] == ')' {
			*currId++
		}
		return "-" + negativeNum
	default:
		token := string(s[*currId])
		*currId++
		return token
	}
}

func calculate(lhs float64, rhs float64, oper string) float64 {
	var result float64

	switch oper {
	case "+":
		result = lhs + rhs
	case "-":
		result = lhs - rhs
	case "*":
		result = lhs * rhs
	case "/":
		result = lhs / rhs
	}

	return result
}

func extractAndCalculate(numSt *Stack[float64], operSt *Stack[string]) float64 {
	rhs, _ := numSt.Pop()
	lhs, _ := numSt.Pop()
	oper, _ := operSt.Pop()
	return calculate(lhs, rhs, oper)
}

// парсинг выражения в стек чисел и стек операторов, вычисление
func ParseExpression(expression string) float64 {
	numStack := NewStack[float64]()
	operStack := NewStack[string]()

	currId := 0
	for currId < len(expression) {
		token := GetNextToken(expression, &currId)
		num, parseErr := strconv.ParseFloat(token, 64)
		switch {
		case parseErr == nil:
			numStack.Push(num)
		case token == "(":
			operStack.Push(token)
		case token == ")":
			lastOper, _ := operStack.Peek()
			for lastOper != "(" {
				numStack.Push(extractAndCalculate(&numStack, &operStack))
				lastOper, _ = operStack.Peek()
			}
			_, _ = operStack.Pop()
		case isOperator[string](token):
			if operStack.isEmpty() {
				operStack.Push(token)
				break
			}
			lastOper, _ := operStack.Peek()
			if lastOper == "(" {
				operStack.Push(token)
			} else {
				if operPriority[token] > operPriority[lastOper] {
					operStack.Push(token)
				} else {
					numStack.Push(extractAndCalculate(&numStack, &operStack))
					operStack.Push(token)
				}
			}
		}
	}

	for !operStack.isEmpty() {
		numStack.Push(extractAndCalculate(&numStack, &operStack))
	}
	res, _ := numStack.Pop()
	return res
}

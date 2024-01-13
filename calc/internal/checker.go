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
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type BaseChecker struct {
}

type InvalidBracketExpression struct {
}

type InvalidSymbol struct {
	sym rune
}

type TooShortExpression struct {
}

type ZeroDivisionError struct {
}

func (e *TooShortExpression) Error() string {
	return "Too short expression"
}

func (e *InvalidSymbol) Error() string {
	return fmt.Sprintf("Invalid symbol: \"%s\"\n", string(e.sym))
}

func (e *InvalidBracketExpression) Error() string {
	return "Invalid bracket sequence"
}

func (t *ZeroDivisionError) Error() string {
	return "Zero division detected"
}

type ParseError struct {
}

func (pe *ParseError) Error() string {
	return ""
}

func isOperator[T rune | string | byte](sym T) bool {
	return strings.Contains("+-*/", string(sym))
}

func hasBrackets(s string) bool {
	for _, sym := range s {
		if sym == '(' || sym == ')' {
			return true
		}
	}

	return false
}

func isClosing(s string) bool {
	st := NewStack[rune]()

	for _, sym := range s {
		if strings.ContainsAny(string(sym), "[]{}") {
			return false
		} else if sym == '(' {
			st.Push(sym)
		} else if sym == ')' {
			if bracket, err := st.Pop(); err != nil || bracket != '(' {
				return false
			}
		}
	}

	if st.isEmpty() {
		return true
	} else {
		return false
	}
}

func CheckBetweenBrackets(s string) bool {
	for i := len(s) - 1; i >= 0; i-- {

		if s[i] == ')' {
			if isOperator(s[i-1]) {
				return false
			}

			var hasNum, hasOperator bool
			for j := i - 1; s[j] != '('; j-- {
				if unicode.IsDigit(rune(s[j])) {
					hasNum = true
				} else if isOperator(s[j]) {
					hasOperator = true
				}
			}
			if !hasNum || !hasOperator {
				return false
			}
		}
	}

	return true
}

func CheckBrackets(s string) error {

	if hasBrackets(s) {
		if !isClosing(s) || !CheckBetweenBrackets(s) {
			return &InvalidBracketExpression{}
		} else {
			return nil
		}
	}

	return nil
}

func CheckOperator(s string, operId int) bool {
	if operId == 0 ||
		operId == len(s)-1 ||
		isOperator(s[operId-1]) || // ++
		isOperator(s[operId+1]) {
		return false
	}

	if !unicode.IsDigit(rune(s[operId-1])) && !unicode.IsDigit(rune(s[operId+1])) {
		if !(s[operId-1] == ')' && s[operId+1] == '(' ||
			s[operId] == '-' && s[operId-1] == '(' && s[operId+1] == '.' ||
			s[operId-1] == ')' && s[operId+1] == '.' ||
			s[operId-1] == '.' && s[operId+1] == '(') { // (2+3)(3-1); 2+(-.5); 3*.5+1; 10.+(-5)
			return false
		}
	} else {
		if s[operId] == '/' && s[operId+1] == '0' { // {.Num}/0
			return false
		}
		if (s[operId-1] == '(' && s[operId] != '-') || s[operId+1] == ')' { // (+...; ...-)
			return false
		}
	}

	return true
}

func CheckDotNeighbours(s string, dotId int) bool {
	if dotId == len(s)-1 {
		if !unicode.IsDigit(rune(s[dotId-1])) {
			return false
		}
	} else {
		if dotId == 0 {
			if !unicode.IsDigit(rune(s[dotId+1])) {
				return false
			}
		} else if !unicode.IsDigit(rune(s[dotId-1])) && !unicode.IsDigit(rune(s[dotId+1])) {
			return false
		}
	}

	return true
}

func CheckDot(s string) bool {
	for id, sym := range s {
		if sym == '.' {
			if !CheckDotNeighbours(s, id) {
				return false
			}

			for i := id + 1; i != len(s) && !isOperator(s[i]); i++ {
				if s[i] == '.' {
					return false
				}
			}
		}
	}
	return true
}

func CheckAllowedSymbols(s string) (rune, bool) {
	for _, sym := range s {
		if !strings.Contains("0123456789+-*/().", string(sym)) {
			return sym, false
		}
	}

	return *new(rune), true
}

func CheckLength(exp string) bool {
	numCount := 0
	var hasOperator bool
	for i := 0; i < len(exp); {
		if isOperator(exp[i]) {
			hasOperator = true
			i++
			continue
		}

		var tmpNum string
		for ; i < len(exp) && unicode.IsDigit(rune(exp[i])); i++ {
			tmpNum += string(exp[i])
		}
		if tmpNum != "" {
			if _, err := strconv.ParseFloat(tmpNum, 64); err == nil {
				numCount++
			}
		} else {
			i++
		}
	}
	if numCount < 2 || !hasOperator {
		return false
	} else {
		return true
	}
}

/*
проверка на:
1. отсутствие недопустимых символов
2. невалидные скобочные последовательности и невалидное содержимое внутри них
3. длину выражения
4. наличие и расположение операторов, плавающей точки
*/
func (ch *BaseChecker) Check(exp string) error {
	if wrongSym, ok := CheckAllowedSymbols(exp); !ok {
		return &InvalidSymbol{sym: wrongSym}
	}
	if err := CheckBrackets(exp); err != nil {
		return err
	}
	if !CheckLength(exp) {
		return &TooShortExpression{}
	}
	if !CheckDot(exp) {
		return &InvalidSymbol{sym: '.'}
	}

	for id, sym := range exp {
		if isOperator(sym) && !CheckOperator(exp, id) {
			return &InvalidSymbol{sym: sym}
		}
	}

	return nil
}

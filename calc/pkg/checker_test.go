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
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsClosing(t *testing.T) {
	cases := []struct {
		s        string
		expected bool
	}{
		{
			s:        "(((())))",
			expected: true,
		},
		{
			s:        "()",
			expected: true,
		},
		{
			s:        "()()()",
			expected: true,
		},
		{
			s:        "(())(())(((())))",
			expected: true,
		},
		{
			s:        "(1+2)-3",
			expected: true,
		},
		{
			s:        "((1+(2+3))-3)+((1+3)*2)",
			expected: true,
		},
		{
			s:        "without_brackets",
			expected: true,
		},
		{
			s:        "",
			expected: true,
		},
		{
			s:        "((-1)+(1+2))+3*((2+(3+1)))",
			expected: true,
		},
		{
			s:        "[2+3]*{3-1}",
			expected: false,
		},
		{
			s:        "(((",
			expected: false,
		},
		{
			s:        "(((((()))",
			expected: false,
		},
		{
			s:        ")))",
			expected: false,
		},
		{
			s:        "()()((())",
			expected: false,
		},
	}

	for _, tCase := range cases {
		t.Run(fmt.Sprintf("Testing isClosing: %+v", tCase), func(t *testing.T) {
			require.Equal(t, tCase.expected, isClosing(tCase.s))
		})
	}
}

func TestCheckBetweenBrackets(t *testing.T) {
	cases := []struct {
		exp      string
		expected bool
	}{
		{
			exp:      "(1+2)+3",
			expected: true,
		},
		{
			exp:      "",
			expected: true,
		},
		{
			exp:      "((-1)+(1+2))+3*((2+(3+1)))",
			expected: true,
		},
		{
			exp:      "((((2 + (((3+4)))))))",
			expected: true,
		},
		{
			exp:      "without brackets",
			expected: true,
		},
		{
			exp:      "2*((((3+1))))+(1)",
			expected: false,
		},
		{
			exp:      "((-1)+(1+2))+3*((2+(3+1)*(2*(5*7))))+((((1+))))",
			expected: false,
		},
		{
			exp:      "()+2",
			expected: false,
		},
		{
			exp:      "(1)+2+3",
			expected: false,
		},
		{
			exp:      "2*(((())))",
			expected: false,
		},
		{
			exp:      "2*((((3+1))))+()",
			expected: false,
		},
		{
			exp:      "2*((((3+1))))+(+)",
			expected: false,
		},
		{
			exp:      "(2+3)+(-2)((2+)*(()))",
			expected: false,
		},
		{
			exp:      "(not_opertor_or_digit)",
			expected: false,
		},
	}

	for _, tCase := range cases {
		t.Run(fmt.Sprintf("Testing CheckBetweenBrackets, %+v", tCase), func(t *testing.T) {
			require.Equal(t, tCase.expected, CheckBetweenBrackets(tCase.exp))
		})
	}
}

func TestCheckOperator(t *testing.T) {

	cases := []struct {
		exp      string
		id       int
		expected bool
	}{
		{
			exp:      "1+2+3",
			id:       1,
			expected: true,
		},
		{
			exp:      "1+2+3",
			id:       3,
			expected: true,
		},
		{
			exp:      "(1+2)+3",
			id:       2,
			expected: true,
		},
		{
			exp:      "(1+2)+3",
			id:       5,
			expected: true,
		},
		{
			exp:      "((-1)+(1+2))+3*((2+(3+1)))",
			id:       2,
			expected: true,
		},
		{
			exp:      "((-1)+(1+2))+3*((2+(3+1)))",
			id:       5,
			expected: true,
		},
		{
			exp:      "((-1)+(1+2))+3*((2+(3+1)))",
			id:       8,
			expected: true,
		},
		{
			exp:      "((-1)+(1+2))+3*((2+(3+1)))",
			id:       12,
			expected: true,
		},
		{
			exp:      "((-1)+(1+2))+3*((2+(3+1)))",
			id:       14,
			expected: true,
		},
		{
			exp:      "((-1)+(1+2))+3*((2++(3+1)))",
			id:       18,
			expected: false,
		},
		{
			exp:      "((-1)+(1+2))+3*((2++(3+1)))",
			id:       18,
			expected: false,
		},
		{
			exp:      "((-1)+(1+2))+3*((2++(3+1)))",
			id:       22,
			expected: true,
		},
		{
			exp:      "((-1)+(1+2))+3*((2++(3+1*)))",
			id:       24,
			expected: false,
		},
		{
			exp:      "+2*3-(-4)+2-",
			id:       0,
			expected: false,
		},
		{
			exp:      "+2*3-(-4)+2-",
			id:       2,
			expected: true,
		},
		{
			exp:      "+2*3-(-4)+2-",
			id:       4,
			expected: true,
		},
		{
			exp:      "+2*3-(-4)+2-",
			id:       6,
			expected: true,
		},
		{
			exp:      "+2*3-(-4)+2-",
			id:       9,
			expected: true,
		},
		{
			exp:      "+2*3-(-4)+2-",
			id:       11,
			expected: false,
		},
		{
			exp:      "+2*3-(-4)+2-(3+4-)",
			id:       11,
			expected: true,
		},
		{
			exp:      "+2*3-(-4)+2-(3+4-)",
			id:       14,
			expected: true,
		},
		{
			exp:      "+2*3-(-4)+2-(3+4-)",
			id:       16,
			expected: false,
		},
	}

	for _, tCase := range cases {
		t.Run(fmt.Sprintf("Testing CheckOperator, %+v", tCase), func(t *testing.T) {
			require.Equal(t, tCase.expected, CheckOperator(tCase.exp, tCase.id))
		})
	}
}

func TestCheckDot(t *testing.T) {
	cases := []struct {
		exp      string
		expected bool
	}{
		{
			exp:      "1.2+3.5+4.",
			expected: true,
		},
		{
			exp:      "1.2+3.5+4.",
			expected: true,
		},
		{
			exp:      "1.2+3.5+4.",
			expected: true,
		},
		{
			exp:      "1.2+3.5+4.",
			expected: true,
		},
		{
			exp:      "1.2+3.5+4.",
			expected: true,
		},
		{
			exp:      ".21*(.5+4.5)/5.",
			expected: true,
		},
		{
			exp:      ".21*(.5+4.5)/5.",
			expected: true,
		},
		{
			exp:      ".21*(.5+4.5)/5.",
			expected: true,
		},
		{
			exp:      ".21*(.5+4.5)/5.",
			expected: true,
		},
		{
			exp:      "123.2222+2.232.",
			expected: false,
		},
		{
			exp:      ".+3..*2",
			expected: false,
		},
		{
			exp:      ".+3..*2",
			expected: false,
		},
		{
			exp:      "(.)",
			expected: false,
		},
		{
			exp:      "(.2.3+1)",
			expected: false,
		},
		{
			exp:      "(.2121221.3213321.+1)",
			expected: false,
		},
		{
			exp:      "123.2222.35+2",
			expected: false,
		},
	}

	for _, tCase := range cases {
		t.Run(fmt.Sprintf("Testing CheckDot: %+v", tCase), func(t *testing.T) {
			require.Equal(t, tCase.expected, CheckDot(tCase.exp))
		})
	}
}

func TestCheckLength(t *testing.T) {
	cases := []struct {
		expression string
		expected   bool
	}{
		{
			expression: "1+2",
			expected:   true,
		},
		{
			expression: "11231-212",
			expected:   true,
		},
		{
			expression: "1",
			expected:   false,
		},
		{
			expression: "11231212",
			expected:   false,
		},
		{
			expression: "11231212+",
			expected:   false,
		},
		{
			expression: "+11231212",
			expected:   false,
		},
		{
			expression: "",
			expected:   false,
		},
		{
			expression: "+",
			expected:   false,
		},
		{
			expression: ".",
			expected:   false,
		},
	}

	for _, tCase := range cases {
		t.Run(fmt.Sprintf("Testing CheckLength: %+v", tCase), func(t *testing.T) {
			require.Equal(t, tCase.expected, CheckLength(tCase.expression))
		})
	}
}

func TestCheck(t *testing.T) {
	cases := []struct {
		expression  string
		expectedErr error
	}{
		{
			expression:  "1+2+3",
			expectedErr: nil,
		},
		{
			expression:  "1+(2+3)",
			expectedErr: nil,
		},
		{
			expression:  "(1+2+3)",
			expectedErr: nil,
		},
		{
			expression:  "((1+2+3))",
			expectedErr: nil,
		},
		{
			expression:  "(1+2)*((3+4)+(3/(1+(4-2))))",
			expectedErr: nil,
		},
		{
			expression:  "(1+2)((3+4)(3/(1+(4-2))))",
			expectedErr: nil,
		},
		{
			expression:  "(-1)+2+3",
			expectedErr: nil,
		},
		{
			expression:  "(-1+3)+2+3",
			expectedErr: nil,
		},
		{
			expression:  "(-1+3)+2+3",
			expectedErr: nil,
		},
		{
			expression:  "(-1.5+3.25)+2+3",
			expectedErr: nil,
		},
		{
			expression:  ".225+(-12.)*3",
			expectedErr: nil,
		},
		{
			expression:  "2+3*.5",
			expectedErr: nil,
		},
		{
			expression:  "(3+.23+5)*3",
			expectedErr: nil,
		},
		{
			expression:  "(3.+23+5)*3",
			expectedErr: nil,
		},
		{
			expression:  "",
			expectedErr: &TooShortExpression{},
		},
		{
			expression:  ".",
			expectedErr: &TooShortExpression{},
		},
		{
			expression:  "+",
			expectedErr: &TooShortExpression{},
		},
		{
			expression:  ".+0",
			expectedErr: &TooShortExpression{},
		},
		{
			expression:  "2+",
			expectedErr: &TooShortExpression{},
		},
		{
			expression:  ".",
			expectedErr: &TooShortExpression{},
		},
		{
			expression:  "21231321.2231.21",
			expectedErr: &TooShortExpression{},
		},
		{
			expression:  "-1+3+2+3",
			expectedErr: &InvalidSymbol{sym: '-'},
		},
		{
			expression:  "((-1)+3*)+2+(-3)(-3)*2",
			expectedErr: &InvalidBracketExpression{},
		},
		{
			expression:  "()",
			expectedErr: &InvalidBracketExpression{},
		},
		{
			expression:  "(1)+2+3",
			expectedErr: &InvalidBracketExpression{},
		},
		{
			expression:  "(1+)+2+3",
			expectedErr: &InvalidBracketExpression{},
		},
		{
			expression:  "()+2+3",
			expectedErr: &InvalidBracketExpression{},
		},
		{
			expression:  "((()))+2+3",
			expectedErr: &InvalidBracketExpression{},
		},
		{
			expression:  "(1+3*)+2+3",
			expectedErr: &InvalidBracketExpression{},
		},
		{
			expression:  "((1+3)+2+3",
			expectedErr: &InvalidBracketExpression{},
		},
		{
			expression:  "))(1+3)+2+3",
			expectedErr: &InvalidBracketExpression{},
		},
		{
			expression:  "(1+3)+2+3)",
			expectedErr: &InvalidBracketExpression{},
		},
		{
			expression:  "(1+3)+(2+3",
			expectedErr: &InvalidBracketExpression{},
		},
		{
			expression:  "(1+3)+(2++3",
			expectedErr: &InvalidBracketExpression{},
		},
		{
			expression:  "(1+3)+(2+323.3.3",
			expectedErr: &InvalidBracketExpression{},
		},
		{
			expression:  ".(1+3)+(2+323",
			expectedErr: &InvalidBracketExpression{},
		},
		{
			expression:  "(1+3)+2+3(",
			expectedErr: &InvalidBracketExpression{},
		},
		{
			expression:  "(1+3)+2+3()",
			expectedErr: &InvalidBracketExpression{},
		},
		{
			expression:  "(3+5.23/)*3",
			expectedErr: &InvalidBracketExpression{},
		},
		{
			expression:  "25.3+4/0",
			expectedErr: &InvalidSymbol{'/'},
		},
		{
			expression:  "2..+3*5",
			expectedErr: &InvalidSymbol{'.'},
		},
		{
			expression:  ".+3*5",
			expectedErr: &InvalidSymbol{'.'},
		},
		{
			expression:  "2.3+3*.",
			expectedErr: &InvalidSymbol{'.'},
		},
		{
			expression:  "333+.+3*5",
			expectedErr: &InvalidSymbol{'.'},
		},
		{
			expression:  "222.3.5+123",
			expectedErr: &InvalidSymbol{'.'},
		},
		{
			expression:  "(3++5)*3",
			expectedErr: &InvalidSymbol{'+'},
		},
		{
			expression:  "(3+5)*-3",
			expectedErr: &InvalidSymbol{'*'},
		},
		{
			expression:  "+(3+5)*3",
			expectedErr: &InvalidSymbol{'+'},
		},
		{
			expression:  "(3+5)*3-",
			expectedErr: &InvalidSymbol{'-'},
		},
		{
			expression:  "(+3+5)*3",
			expectedErr: &InvalidSymbol{'+'},
		},
		{
			expression:  "+3+5*3",
			expectedErr: &InvalidSymbol{'+'},
		},
	}

	checker := BaseChecker{}
	for _, tCase := range cases {
		t.Run(fmt.Sprintf("Testing Check: %s", tCase.expression), func(t *testing.T) {
			require.Equal(t, tCase.expectedErr, checker.Check(tCase.expression))
		})
	}
}

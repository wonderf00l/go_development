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
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func makeMockReader(s string) io.Reader {
	stringsReader := strings.NewReader(s)
	return stringsReader
}

func makeMockWriter() *bytes.Buffer {
	return new(bytes.Buffer)
}

func NewInt(num int) *int {
	var i int = num
	return &i
}

func TestGetNextToken(t *testing.T) {
	cases := []struct {
		s     string
		id    *int
		newId int
		token string
	}{
		{
			s:     "(2+3)*1.25-(-1234.543)",
			id:    NewInt(0),
			newId: 1,
			token: "(",
		},
		{
			s:     "(2+3)*1.25-(-1234.543)",
			id:    NewInt(1),
			newId: 2,
			token: "2",
		},
		{
			s:     "(2+3)*1.25-(-1234.543)",
			id:    NewInt(2),
			newId: 3,
			token: "+",
		},
		{
			s:     "(2+3)*1.25-(-1234.543)",
			id:    NewInt(6),
			newId: 10,
			token: "1.25",
		},
		{
			s:     "(2+3)*1.25-(-1234.543)",
			id:    NewInt(11),
			newId: len("(2+3)*1.25-(-1234.543)"),
			token: "-1234.543",
		},
		{
			s:     "(-.4321+23)",
			id:    NewInt(0),
			newId: 7,
			token: "-.4321",
		},
		{
			s:     "((-.4321)+23)",
			id:    NewInt(1),
			newId: 9,
			token: "-.4321",
		},
		{
			s:     "6*(3+((-1)-3)*(-1))/2-(.5+1.5)",
			id:    NewInt(6),
			newId: 10,
			token: "-1",
		},
		{
			s:     "(-.4321+23)",
			id:    NewInt(1234),
			newId: 1234,
			token: "",
		},
		{
			s:     "-.4321+23",
			id:    NewInt(0),
			newId: 1,
			token: "-",
		},
	}

	for _, tCase := range cases {
		t.Run(fmt.Sprintln("Testing GetNextToken: ", tCase.s, *tCase.id), func(t *testing.T) {
			actualToken := GetNextToken(tCase.s, tCase.id)
			require.Equal(t, tCase.token, actualToken)
			require.Equal(t, tCase.newId, *tCase.id)
		})
	}
}

func TestParseExpression(t *testing.T) {
	cases := []struct {
		expression string
		result     float64
	}{
		{
			expression: "1+2+3",
			result:     float64(6),
		},
		{
			expression: "1+(2-3)",
			result:     float64(0),
		},
		{
			expression: "1+(2-3)*5",
			result:     float64(-4),
		},
		{
			expression: "6*(3+((-1)-3)*(-1))/2-(.5+1.5)",
			result:     float64(19),
		},
		{
			expression: "6.455*(3+((-1)-3)*(-1))/2-(.5+1.5)",
			result:     float64(20.5925),
		},
		{
			expression: "((6*(3+((-1)-3)*(-1))/2-(.5+1.5)))",
			result:     float64(19),
		},
		{
			expression: "((6*(3+((-1)-3)*(-1))/2-(.5+1.5)))*(1+1)",
			result:     float64(38),
		},
		{
			expression: "2+(-.5)*.5/1",
			result:     float64(1.75),
		},
		{
			expression: "225.0+(-.5)*(-1)-10.+(-5)",
			result:     float64(210.5),
		},
	}

	for _, tCase := range cases {
		t.Run(fmt.Sprintln("Testing ParseExpression: ", tCase.expression), func(t *testing.T) {
			require.Equal(t, tCase.result, ParseExpression(tCase.expression))
		})
	}
}

func TestCalculator(t *testing.T) {
	config := NewConfig()
	cases := []struct {
		in       io.Reader
		out      *bytes.Buffer
		expected string
		Err      error
	}{
		{
			in:       makeMockReader("1+2+3"),
			out:      makeMockWriter(),
			expected: "6\n",
			Err:      nil,
		},
		{
			in:       makeMockReader("1+(2-3)"),
			out:      makeMockWriter(),
			expected: "0\n",
			Err:      nil,
		},
		{
			in:       makeMockReader("1+(2-3)*5"),
			out:      makeMockWriter(),
			expected: "-4\n",
			Err:      nil,
		},
		{
			in:       makeMockReader("6*(3+((-1)-3)*(-1))/2-(.5+1.5)"),
			out:      makeMockWriter(),
			expected: "19\n",
			Err:      nil,
		},
		{
			in:       makeMockReader("6.455*(3+((-1)-3)*(-1))/2-(.5+1.5)"),
			out:      makeMockWriter(),
			expected: "20.5925\n",
			Err:      nil,
		},
		{
			in:       makeMockReader("((6*(3+((-1)-3)*(-1))/2-(.5+1.5)))"),
			out:      makeMockWriter(),
			expected: "19\n",
			Err:      nil,
		},
		{
			in:       makeMockReader("((6*(3+((-1)-3)*(-1))/2-(.5+1.5)))*(1+1)"),
			out:      makeMockWriter(),
			expected: "38\n",
			Err:      nil,
		},
		{
			in:       makeMockReader("2+(-.5)*.5/1"),
			out:      makeMockWriter(),
			expected: "1.75\n",
			Err:      nil,
		},
		{
			in:       makeMockReader("225.0+(-.5)*(-1)-10.+(-5)"),
			out:      makeMockWriter(),
			expected: "210.5\n",
			Err:      nil,
		},
		{
			in:       makeMockReader(""),
			out:      makeMockWriter(),
			expected: "",
			Err:      nil,
		},
		{
			in:       makeMockReader("22/0"),
			out:      makeMockWriter(),
			expected: "",
			Err:      &InvalidSymbol{'/'},
		},
		{
			in:       makeMockReader("22++123"),
			out:      makeMockWriter(),
			expected: "",
			Err:      &InvalidSymbol{'+'},
		},
		{
			in:       makeMockReader("*22-123"),
			out:      makeMockWriter(),
			expected: "",
			Err:      &InvalidSymbol{'*'},
		},
		{
			in:       makeMockReader("22+123/"),
			out:      makeMockWriter(),
			expected: "",
			Err:      &InvalidSymbol{'/'},
		},
		{
			in:       makeMockReader("22.11.12+123"),
			out:      makeMockWriter(),
			expected: "",
			Err:      &InvalidSymbol{'.'},
		},
		{
			in:       makeMockReader("22..12+123"),
			out:      makeMockWriter(),
			expected: "",
			Err:      &InvalidSymbol{'.'},
		},
		{
			in:       makeMockReader("(+22.12*3/2)+123"),
			out:      makeMockWriter(),
			expected: "",
			Err:      &InvalidSymbol{'+'},
		},
		{
			in:       makeMockReader(" "),
			out:      makeMockWriter(),
			expected: "",
			Err:      &InvalidSymbol{' '},
		},
		{
			in:       makeMockReader("23.y+4"),
			out:      makeMockWriter(),
			expected: "",
			Err:      &InvalidSymbol{'y'},
		},
		{
			in:       makeMockReader("\n"),
			out:      makeMockWriter(),
			expected: "",
			Err:      &TooShortExpression{},
		},
		{
			in:       makeMockReader(".0"),
			out:      makeMockWriter(),
			expected: "",
			Err:      &TooShortExpression{},
		},
		{
			in:       makeMockReader("(1+)"),
			out:      makeMockWriter(),
			expected: "",
			Err:      &InvalidBracketExpression{},
		},
		{
			in:       makeMockReader("((1+2)"),
			out:      makeMockWriter(),
			expected: "",
			Err:      &InvalidBracketExpression{},
		},
		{
			in:       makeMockReader("1+3+4)"),
			out:      makeMockWriter(),
			expected: "",
			Err:      &InvalidBracketExpression{},
		},
		{
			in:       makeMockReader("(12312)+33"),
			out:      makeMockWriter(),
			expected: "",
			Err:      &InvalidBracketExpression{},
		},
		{
			in:       makeMockReader("()*(22+3)"),
			out:      makeMockWriter(),
			expected: "",
			Err:      &InvalidBracketExpression{},
		},
		{
			in:       makeMockReader("(()()))"),
			out:      makeMockWriter(),
			expected: "",
			Err:      &InvalidBracketExpression{},
		},
		{
			in:       makeMockReader("(1+(2+(1+()))"),
			out:      makeMockWriter(),
			expected: "",
			Err:      &InvalidBracketExpression{},
		},
	}

	for _, tCase := range cases {
		t.Run(fmt.Sprintf("Testing Calculator: %+v\n", tCase.in), func(t *testing.T) {
			err := Calculator(tCase.in, tCase.out, config)
			require.Equal(t, tCase.Err, err)
			require.Equal(t, tCase.expected, tCase.out.String())
		})
	}
}

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

func TestStack(t *testing.T) {
	intSt, stringSt := NewStack[int](), NewStack[string]()
	expInt, expString := make([]int, 0), make([]string, 0)

	for i := 0; i < 10; i++ {
		intSt.Push(i)
		expInt = append(expInt, i)
		t.Run(fmt.Sprint("Testing IntStack: ", intSt), func(t *testing.T) {
			require.Equal(t, expInt, []int(intSt))

			val, err := intSt.Peek()
			require.Equal(t, nil, err)
			require.Equal(t, expInt[len(expInt)-1], val)
			require.False(t, intSt.isEmpty())
		})
	}

	for i := 0; i < 10; i++ {
		t.Run(fmt.Sprint("Testing IntStack on Popping: ", intSt), func(t *testing.T) {
			val, err := intSt.Pop()
			require.Equal(t, expInt[len(expInt)-1], val)

			expInt = expInt[:len(expInt)-1]
			require.Equal(t, expInt, []int(intSt))
			require.Equal(t, nil, err)
		})
	}
	require.True(t, intSt.isEmpty())

	strs := []string{"input", "variable", "into", "stack"}

	for _, str := range strs {
		stringSt.Push(str)
		expString = append(expString, str)
		t.Run(fmt.Sprint("Testing StringStack: ", stringSt), func(t *testing.T) {
			require.Equal(t, expString, []string(stringSt))

			val, err := stringSt.Peek()
			require.Equal(t, nil, err)
			require.Equal(t, expString[len(expString)-1], val)
			require.False(t, stringSt.isEmpty())
		})
	}

	for i := len(expString) - 1; i >= 0; i-- {
		t.Run(fmt.Sprint("Testing StringStack on Popping: ", stringSt), func(t *testing.T) {
			val, err := stringSt.Pop()
			require.Equal(t, expString[len(expString)-1], val)

			expString = expString[:len(expString)-1]
			require.Equal(t, expString, []string(stringSt))
			require.Equal(t, nil, err)
		})
	}
	require.True(t, stringSt.isEmpty())

	val, err := intSt.Pop()
	require.Equal(t, &EmptyStackError{}, err)
	require.Equal(t, 0, val)

	val_, err := stringSt.Pop()
	require.Equal(t, &EmptyStackError{}, err)
	require.Equal(t, "", val_)
}

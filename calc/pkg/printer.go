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
	"io"
)

type BasePrinter struct {
}

func (p *BasePrinter) Print(num float64, out io.Writer) error {
	_, err := fmt.Fprintln(out, num)
	return err
}

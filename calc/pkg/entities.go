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
	"io"
	"os"
)

type Checker interface {
	Check(string) error
}

type Printer interface {
	Print(float64, io.Writer) error
}

type Config struct {
	InputSteam   io.Reader
	OutputStream io.Writer
	Checker      Checker
	Printer      Printer
}

func NewConfig() *Config {
	return &Config{
		InputSteam:   os.Stdin,
		OutputStream: os.Stdout,
		Checker:      &BaseChecker{},
		Printer:      &BasePrinter{},
	}
}

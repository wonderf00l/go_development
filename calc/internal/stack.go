/*
Пакет calc содержит:
1. stack.go/stack_test.go - реализаия стека и тесты
2. entities.go - основные сущности приложения
3. checker.go/checker_test.go - реализация методов, проверяющих входное выражение
4. printer.go - реализация вывода результата
5. calc.go/calc_test.go - основные функции для чтения выражения, парсинга и вычисления результата
*/
package calc

type Stack[T comparable] []T

type EmptyStackError struct {
}

func (e *EmptyStackError) Error() string {
	return "Trying to delete value from an empty stack"
}

func NewStack[T comparable]() Stack[T] {
	return make(Stack[T], 0)
}

func (st *Stack[T]) Push(val T) {
	*st = append(*st, val)
}

func (st *Stack[T]) Pop() (T, error) {
	if len(*st) == 0 {
		return *new(T), &EmptyStackError{}
	}

	popped := (*st)[len(*st)-1]
	*st = (*st)[:len(*st)-1]
	return popped, nil
}

func (st Stack[T]) Peek() (T, error) {
	if len(st) == 0 {
		return *new(T), &EmptyStackError{}
	}

	return st[len(st)-1], nil
}

func (st Stack[T]) isEmpty() bool {
	return len(st) == 0
}

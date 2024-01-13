// Пакет uniq содержит код утилиты, может быть переиспользован в других пакетах
// содержит файлы:
// с логикой извлечения и обработки аргументов командной строки --> options.go
// с основными сущностями приложения --> entities.go
// с реализацией методов основных сущностей и их вспомогательными функциями --> methods.go
// с главными функциями приложения, куда импортируются все сущности --> uniquify.go
package uniq

import "fmt"

type ZeroStringError struct{}

func (e *ZeroStringError) Error() string {
	return "Got zero string: \"\""
}

type WrongTemplateError struct {
	S, Substr string
}

func (e *WrongTemplateError) Error() string {
	return fmt.Sprintf("\"%s\" not in \"%s\"", e.Substr, e.S)
}

type InitError struct {
	Data interface{}
}

func (err *InitError) Error() string {
	return fmt.Sprintf("Failed to provide data to Formater: %v", err.Data)
}

type IncompatibleOptionsError struct {
	Opts map[string]string
}

func (inc *IncompatibleOptionsError) Error() string {
	return fmt.Sprintf("Got incompatible arguments: %v", inc.Opts)
}

type ExcessOptionsError struct {
	Opts []string
}

func (exc *ExcessOptionsError) Error() string {
	return fmt.Sprintf("Got excess arguments: %v", exc.Opts)
}

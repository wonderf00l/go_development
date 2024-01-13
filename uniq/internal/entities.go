// Пакет uniq содержит код утилиты, может быть переиспользован в других пакетах
// содержит файлы:
// с логикой извлечения и обработки аргументов командной строки --> options.go
// с основными сущностями приложения --> entities.go
// с реализацией методов основных сущностей и их вспомогательными функциями --> methods.go
// с главными функциями приложения, куда импортируются все сущности --> uniquify.go
package uniq

import "io"

// Comparator сравнивает строки
type Comparator interface {
	Compare(string, string) bool
}

// Processor обрабатывает строку
type Processor interface {
	Process(string) string
}

// Formater форматирует строку
type Formater interface {
	Format(string) string
}

// Printer пишет данные строки в выходной поток
type Printer interface {
	Print(string, io.Writer) error
}

type LexicographicalComparator struct{}

type BaseProcessor struct {
	IgnoreReg   bool
	FieldOffset uint
	CharOffset  uint
}

type BaseFormater struct {
	CountRow       bool
	RowFreq        int
	RepetitiveOnly bool
	UniqOnly       bool
}

type BasePrinter struct{}

type Config struct {
	InputStream  io.ReadCloser
	OutputStream io.WriteCloser
	Comp         Comparator
	Proc         Processor
	Formater     Formater
	Printer      Printer
}

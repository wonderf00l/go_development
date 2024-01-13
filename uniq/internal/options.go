// Пакет uniq содержит код утилиты, может быть переиспользован в других пакетах
// содержит файлы:
// с логикой извлечения и обработки аргументов командной строки --> options.go
// с основными сущностями приложения --> entities.go
// с реализацией методов основных сущностей и их вспомогательными функциями --> methods.go
// с главными функциями приложения, куда импортируются все сущности --> uniquify.go
package uniq

import (
	"flag"
)

// Options хранит опции, поулчаемы из командной строки
type Options struct {
	InputFile      string
	OutputFile     string
	CountRow       bool
	RepetitiveOnly bool
	UniqOnly       bool
	FieldOffset    uint
	CharOffset     uint
	IgnoreReg      bool
}

var incompatibleOpts = map[string][]string{
	"d": {"u"},
	"u": {"d"},
}

// проверка опций на совместимость
func checkCompability(opts *Options) error {
	if opts.RepetitiveOnly && opts.UniqOnly {
		return &IncompatibleOptionsError{map[string]string{"-d": "-u"}}
	}

	return nil
}

// извлечение опций, парсинг флагов
func (opts *Options) parseFlags() error {
	flag.BoolVar(&opts.CountRow, "c", false, "count repetitive rows in input data")
	flag.BoolVar(&opts.RepetitiveOnly, "d", false, "provide repetitive rows only")
	flag.BoolVar(&opts.UniqOnly, "u", false, "provide uniq rows only")
	flag.BoolVar(&opts.IgnoreReg, "i", false, "ignore input data register")
	flag.UintVar(&opts.FieldOffset, "f", 0, "number of fields to skip(separated by whitespaces)")
	flag.UintVar(&opts.CharOffset, "s", 0, "number of chars to skip")
	flag.Parse()

	if length := len(flag.Args()); length <= 2 {
		opts.InputFile, opts.OutputFile = flag.Arg(0), flag.Arg(1)
		return nil
	}

	return &ExcessOptionsError{flag.Args()[2:]}
}

// создание опций с проверкой входных параметров
func NewOptions() (*Options, error) {
	var opts Options

	if parseErr := opts.parseFlags(); parseErr != nil {
		return &Options{}, parseErr
	}
	if compatErr := checkCompability(&opts); compatErr != nil {
		return &Options{}, compatErr
	}

	return &opts, nil
}

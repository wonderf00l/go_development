// Пакет uniq содержит код утилиты, может быть переиспользован в других пакетах
// содержит файлы:
// с логикой извлечения и обработки аргументов командной строки --> options.go
// с основными сущностями приложения --> entities.go
// с реализацией методов основных сущностей и их вспомогательными функциями --> methods.go
// с главными функциями приложения, куда импортируются все сущности --> uniquify.go
package uniq

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/unicode/norm"
)

// создание нового конфига с учетом полученных опций, инициализация основных сущностей
func NewConfig(opts *Options) (*Config, error) {
	var inputStream io.ReadCloser
	var outputStream io.WriteCloser
	var err error

	if opts.InputFile != "" {
		if inputStream, err = os.Open(opts.InputFile); err != nil {
			log.Printf("Configuring input stream for application: %s\n", err)
			return nil, err
		}
	} else {
		inputStream = os.Stdin
	}

	if opts.OutputFile != "" {
		if outputStream, err = os.Open(opts.OutputFile); err != nil {
			log.Printf("Configuring output stream for application: %s\n", err)
			return nil, err
		}
	} else {
		outputStream = os.Stdout
	}

	return &Config{
		InputStream:  inputStream,
		OutputStream: outputStream,
		Comp:         &LexicographicalComparator{},
		Proc: &BaseProcessor{
			IgnoreReg:   opts.IgnoreReg,
			FieldOffset: opts.FieldOffset,
			CharOffset:  opts.CharOffset,
		},
		Formater: &BaseFormater{CountRow: opts.CountRow, RepetitiveOnly: opts.RepetitiveOnly, UniqOnly: opts.UniqOnly},
		Printer:  &BasePrinter{},
	}, nil
}

func (cfg *Config) Close() {
	cfg.InputStream.Close()
	cfg.OutputStream.Close()
}

// Метод базового процессора: нормализует строку, пропускает FieldOffset полей, CharOffset символов
func (pr *BaseProcessor) Process(s string) string {
	s = norm.NFC.String(s)

	if pr.IgnoreReg {
		s = strings.ToLower(s)
	}
	if pr.FieldOffset != 0 {
		s = skipFields(s, pr.FieldOffset)
	}

	if length := uint(utf8.RuneCountInString(s)); length < pr.CharOffset {
		pr.CharOffset = length // если будет ситуация "string"[len("string"):] --> получим "", runtime ошибки не будет, аналогично и с []rune{}[len([]rune):]
	}

	return string([]rune(s)[pr.CharOffset:])
}

// Функция для конвертации в новую строку с учетом параметра FieldOffset
func skipFields(s string, FieldOffset uint) string {
	if s == "" {
		return ""
	}
	if FieldOffset == 0 {
		return s
	}
	if strings.IndexFunc(s, func(c rune) bool { return c != ' ' && c != '\t' && c != '\n' }) == -1 {
		return ""
	}

	startFrom := 0
	fields := strings.Fields(s)
	if FieldOffset > uint(len(fields)) {
		FieldOffset = uint(len(fields))
	}
	lastSkipField := fields[FieldOffset-1]
	startFrom, _ = getNewStrId(s, lastSkipField, count(lastSkipField, fields[:FieldOffset]))
	startFrom++

	return s[startFrom:] // если будет ситуация "string"[len("string"):] --> получим "", runtime ошибки не будет
}

// Вспомогательная функция для поиска позиции подстроки substr в строке s, с которой будет начинаться новая строка, в которой пропущено FieldOffset полей
// offset - позиция поля в части строки, которую нужно откинуть, иначе говоря, какое именно из полей выбрать, если встретиться несколько одинаковых
// сценарий с возвращением ошибки - для переиспользования, хотя функция SkipFields предусматривает обработку краевых занчений
func getNewStrId(s, substr string, offset int) (int, error) {
	if shouldBeProceed, pos, err := checkCredentials(s, substr, offset); !shouldBeProceed {
		return pos, err
	}

	var currFreq, currPos, endPos int
	for currFreq < offset {
		if endPos >= len(s) {
			return len(s) - 1, nil
		}
		currPos = endPos + strings.Index(s[endPos:], substr) // shift str
		endPos = currPos + (len(substr))
		if currPos == 0 {
			if strings.Contains(" \t", string(s[endPos])) {
				currFreq++
			}
		} else if endPos == len(s) {
			if strings.Contains(" \t", string(s[currPos-1])) {
				currFreq++
			}
		} else if strings.Contains(" \t", string(s[currPos-1])) && strings.Contains(" \t", string(s[endPos])) {
			currFreq++
		}
	}

	return currPos + (len(substr) - 1), nil
}

// вспомогательная функция для подсчета частоты строки в слайсе строк, нужно для высчитывания позиции конкретного поля из множества полей строки(если одинаковых полей несколько)
func count(s string, fields []string) int {
	counter := 0

	for _, str := range fields {
		if str == s {
			counter++
		}
	}

	return counter
}

// проверка строки s и подстроки substr, а также offset на краевые случаи
// возвращает апрув на продолжение поиска новой позиции в обрабатываемой строке
func checkCredentials(s, substr string, offset int) (bool, int, error) {
	shouldBeProceed := true

	if offset == 0 {
		return false, 0, nil
	}
	if s == "" || substr == "" {
		return false, 0, &ZeroStringError{}
	}
	if !strings.Contains(s, substr) {
		return false, 0, &WrongTemplateError{Substr: substr, S: s}
	}
	if s == substr {
		return false, len(s) - 1, nil
	}

	return shouldBeProceed, 0, nil
}

// метод дефолтного компаратора
func (lexComp *LexicographicalComparator) Compare(lhs, rhs string) bool {
	return lhs == rhs
}

// метод дефолтного форматера: конвертирует входную строку s с учетом опций приложения
func (baseF *BaseFormater) Format(s string) string {
	var formated string

	if baseF.RowFreq > 1 && baseF.UniqOnly ||
		baseF.RowFreq <= 1 && baseF.RepetitiveOnly {
		return formated
	}
	if baseF.CountRow {
		return fmt.Sprintf("%7d %s", baseF.RowFreq, s)
	}

	return s
}

// функция для инициализации параметров форматера
// для дефолтного форматера - частоты строки во входном потоке
func ProvideDataToFormater(f Formater, data interface{}) error {
	switch f.(type) {
	case *BaseFormater:
		baseF, ok := f.(*BaseFormater)
		if !ok {
			return &InitError{data}
		}
		switch data.(type) {
		case int:
			data_, ok_ := data.(int)
			if !ok_ {
				return &InitError{data_}
			} else {
				baseF.RowFreq = data_
			}
		}
	}

	return nil
}

// метод дефолтного принтера: пишем в out-поток
func (baesPr *BasePrinter) Print(s string, out io.Writer) error {
	_, err := out.Write([]byte(s))
	return err
}

// Пакет uniq содержит код утилиты, может быть переиспользован в других пакетах
// содержит файлы:
// с логикой извлечения и обработки аргументов командной строки --> options.go
// с основными сущностями приложения --> entities.go
// с реализацией методов основных сущностей и их вспомогательными функциями --> methods.go
// с главными функциями приложения, куда импортируются все сущности --> uniquify.go
package uniq

import (
	"bufio"
	"io"
)

// вспомогательная функция для просмотра входного потока до '\n' включительно с сохранением позиции
func peekLine(in *bufio.Reader) ([]byte, error) {
	var line, symbols []byte
	var err error
	i := 0

	for ; err == nil; i++ {
		symbols, err = in.Peek(1 + i)
		if err == io.EOF || symbols[len(symbols)-1] == 10 {
			break
		}
	}
	if err != nil {
		return []byte{}, err
	}
	line = append(line, symbols...)
	return line, err
}

// построчное чтение входного потока данных, в качестве разделителя - '\n':
// для текущей строки ищем дубликаты с учетом переданных опций
// для сохранения состояния входного потока в случае считывания "новой" строки, не дубликата,
// используется PeekLine, если же строка оказалось дубликатом, сдвигаем позицию во входном потоке
// до следующей строки
// далее строка форматируется и выводится в выходной поток
func Uniquify(in io.Reader, out io.Writer, config *Config) error {
	src := bufio.NewReader(in)
	var err error
	var currLine string

	for err == nil {
		currLine, err = src.ReadString('\n')
		lineCount := 1
		var nextLine []byte
		for err == nil {
			nextLine, err = peekLine(src)
			if err == io.EOF ||
				!config.Comp.Compare(config.Proc.Process(currLine), config.Proc.Process(string(nextLine))) {
				break
			}
			lineCount++
			_, _ = src.ReadString('\n') // прочитываем строку, сдвигаем указатель
		}
		ProvideDataToFormater(config.Formater, lineCount)
		formatedOutput := config.Formater.Format(currLine)
		config.Printer.Print(formatedOutput, out)
	}

	return err
}

package main

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"
)

const goroutinesNum = 3

// каналы потокобезопасны, внутри каналов есть мьютексы

func startWorker(wg *sync.WaitGroup, workerNum int, in <-chan string) {
	defer wg.Done()
	for input := range in {
		fmt.Printf(formatWork(workerNum, input))
		runtime.Gosched() // попробуйте закомментировать; шед закидывает воркер в очередь, чтобы какой-то другой смог перехватить таску и работа велась более равномерно
		// без goshed в основном одна гоуритна будет лочиться в ожидании и сразу анлочиться при поялвении данных
	}
	printFinishWork(workerNum)
}

func main() {
	runtime.GOMAXPROCS(0) // попробуйте с 0 (все доступные) и 1
	// особо не влияет, результат не опредлен, все работают по-разному
	// без goshed - GOMAXPROCS дает гоуртин побежать на др процессорах --> работа более равномерная
	worketInput := make(chan string, 3) // попробуйте увеличить размер канала
	wg := &sync.WaitGroup{}
	for i := 0; i < goroutinesNum; i++ {
		wg.Add(1)
		go startWorker(wg, i, worketInput)
	}

	months := []string{"Январь", "Февраль", "Март",
		"Апрель", "Май", "Июнь",
		"Июль", "Август", "Сентябрь",
		"Октябрь", "Ноябрь", "Декабрь",
	}

	for _, monthName := range months {
		worketInput <- monthName
	}
	close(worketInput) // попробуйте закомментировать
	// буферизированный канал тоже необходимо закрывать, не получится считать только данные буфера в range и выйти
	wg.Wait()
	time.Sleep(time.Millisecond)
}

/*
	МОЖЕМ НЕ ЗАКРЫВАТЬ КАНАЛ ПОСЛЕ ЗАПИСИ ВСЕХ ДАННЫХ В НЕГО
	НО ТОГДА ГОРУТИНА ПРОСТО ЗАВИСНЕТ В FOR-RANGE В ПОПЫТКЕ СЧИТАТЬ ЧТО-ТО ИЗ КАНАЛА
	причем дедлока не будет, если мы не будем ждать все завершения всех горутин через wg
	гоуртины просто заблоктруются и утекут, будут ждать впустую
	и потому без закрытия и wg они только отпринтят информацию о получении значения из канала,
	но не о завершении
*/

func formatWork(in int, input string) string {
	return fmt.Sprintln(strings.Repeat("  ", in), "█",
		strings.Repeat("  ", goroutinesNum-in),
		"th", in,
		"recieved", input)
}

func printFinishWork(in int) {
	fmt.Println(strings.Repeat("  ", in), "█",
		strings.Repeat("  ", goroutinesNum-in),
		"===", in,
		"finished")
}

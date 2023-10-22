package main

import (
	"fmt"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(0)

	in := make(chan int, 1)

	go func(out chan<- int) { // out канал - канал только для чтения
		for i := 0; i <= 10; i++ {
			fmt.Println("before", i)
			out <- i

			fmt.Println("after", i)
		}
		// out <- 12
		close(out) // теперь нельзя писать + нельзя переоткрыть, только создать новый
		fmt.Println("generator finish")
	}(in)

	for i := range in { // вычитвание значений из канала, пока там есть значения и он открыт
		/// итерация по каналу и запись в него происходят асинхронно
		// func и main исполняются, возможно, на разных ядрах
		// когда кладем в канал в func, range ждем данные, сразу берет их
		// после вызова close() range дочитывает все данные из буфера до конца и потом словим событие 'канал closed'
		fmt.Println("\tget", i)
	}

	for { // аналог range
		data, ok := <-in
		fmt.Println(data, ok)
		if !ok {
			fmt.Println("channel closed")
			break
		}
	} // читать из закрытого канала можно (пытатсья чиатть), но будем поулчать ok == false
	// fmt.Scanln()
}

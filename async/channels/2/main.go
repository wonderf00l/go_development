package main

import "fmt"

func sender(c chan int) {
	c <- 1 // len 1, cap 3
	c <- 2 // len 2, cap 3
	c <- 3 // len 3, cap 3
	c <- 4 // <- goroutine blocks here
	close(c)
}

func main() {
	c := make(chan int, 3)

	go sender(c)

	fmt.Printf("Length of channel c is %v and capacity of channel c is %v\n", len(c), cap(c))

	// read values from c (blocked here)
	for val := range c { // видимо, посколько в for - range мы блокируемся
		// до тех пор, пока в канал не будет записано значение
		// в момент блокировки управление передается горутине sender
		// она за время своей работы успевает записать в канал все 3 значения(поэтому и видим при первом выводе len(c) == 3)
		// однако на 4-ом она блокируется в ожидании, что-то
		// какая-то из горутин считает 0-ой элемент буфера
		// то есть планировщик активизирует горутину main
		// в for читаем 3 значения, снова блокируемся
		// просыпается sender, дописывает свою 4, закрывает канал и завершает работу
		// main дочитывает последнее значение из уже закрытого канала
		// выпадаем из for
		fmt.Printf("Length of channel c after value '%v' read is %v\n", val, len(c))
	}
}

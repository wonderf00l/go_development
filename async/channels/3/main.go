package main

import (
	"fmt"
	"runtime"
)

func squares(goroutine_id int, c chan int) {
	for i := 0; i < 4; i++ {
		num := <-c
		fmt.Println("got ", num*num, "in ", goroutine_id)
	}
}

func main() {
	fmt.Println("main() started")
	c := make(chan int, 3)
	go squares(0, c)

	fmt.Println("active goroutines", runtime.NumGoroutine()) // 2
	c <- 1
	c <- 2
	c <- 3
	c <- 4 // blocks here, шедулер паркует галвнуб горутину, пробуждает squares, та читает 3 значения и блокируется
	// пробуждается main, дописывает 4, отходит в ready to run, в running переходит squares и доделывает свою работу
	// она дочиытвает последнее значение из канала и ее цикл заканчивает -> она завершает работу

	fmt.Println("active goroutines", runtime.NumGoroutine()) // 1

	go squares(1, c)

	fmt.Println("active goroutines", runtime.NumGoroutine()) // 2

	c <- 5
	c <- 6
	c <- 7
	c <- 8 // blocks here

	fmt.Println("active goroutines", runtime.NumGoroutine()) // 1
	fmt.Println("main() stopped")
}

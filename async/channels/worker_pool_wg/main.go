package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	in, out := make(chan int, 8), make(chan int, 8)

	wg := &sync.WaitGroup{}
	wg_ := sync.WaitGroup{}
	wg_.Add(1) // функции, принимающие Pointer reciever, могут менять свою структуру, но сама структура, а не указатель на нее, можем вызывать pReceiver методы

	for i := 0; i != 5; i++ {
		wg.Add(1)
		go func(in, out chan int, id int) {
			defer wg.Done()
			fmt.Println("Inside worker №", id)
			// time.Sleep(time.Second)
			for data := range in {
				out <- data * data
				time.Sleep(10 * time.Millisecond) // если не заблокаемся, то какая-то одна горутина быстрее перехватит инициативу и выполнит всю работу
				fmt.Printf("worker %d wrote %d in out\n", id, data*data)
				if data/2 == 4 {
					close(out) // if no close --> deadlock
				}
			}
		}(in, out, i)
	}

	for i := 0; i != 5; i++ {
		in <- i * 2
	}

	close(in) // if no close -> deadlock

	for outputData := range out {
		fmt.Printf("got %d\n", outputData)
	}

	wg.Wait()                                       // т.к. ждем всех, то в итоге все горутины 100 проц отработают минимум до range по каналу, то есть выведут свои принты
	fmt.Println("waiting is finished, terminating") // т.к. in зыкравется, они закончат работу без дедлока
}

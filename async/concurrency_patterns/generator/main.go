package main

import "fmt"

// fib returns a channel which transports fibonacci numbers
func fib(length int) <-chan int {
	// make buffered channel
	c := make(chan int) // конкурентное выполнение, если 0 <= ch_buf_len < length

	// run generation concurrently
	go func() {
		for i, j := 0, 1; i < length; i, j = i+j, i {
			fmt.Printf("goroutine is going to send %d value to channel\n", i)
			c <- i // скорее всего, тут происходит запись значения в структуру ресивера и повторная планировка ресивера, то есть отправка его в состояние ready to run, за это время горутина успевает отпринтиться
			fmt.Printf("goroutine sent %d value to channel\n", i)
		}
		close(c)
	}()

	// return channel
	return c
}

func main() {
	// read 10 fibonacci numbers from channel returned by `fib` function
	for fn := range fib(10) {
		fmt.Println("Current fibonacci number is", fn)
	}
}

// запись в канал в анонимной горутине -> чтение в main горутине -> и так по очереди конкурентно

package main

import "fmt"

func Turnout(inA, inB <-chan int, outA, outB chan<- int) {

	for {
		var (
			data int
			more bool
		)
		select { // т.к. select без default блокирующий
		// то будем ждать появления данных в любом из каналов
		// но если хоть один из каналов будет закрыт
		// и отправит more == false, то сразу выйдем
		// то есть можем не счиать данные с еще не закрытого канала
		case data, more = <-inA:
		case data, more = <-inB:
		}
		if !more {
			return
		}
		select {
		case outA <- data:
		case outB <- data:
		}
	}

}

func main() {
	ch := make(chan int, 10)

	for i := 0; i != 5; i++ {
		ch <- i
	}

	close(ch)

	for data := range ch { // можем дочитать информацию из buffered канала после закрытия
		fmt.Printf("got data %d from channel\n", data)
	}

	unBuf := make(chan int)

	go func(unbuf chan int) {
		unbuf <- 555
		close(unbuf)
	}(unBuf) // закрытие канала лишь означает, что в него больше нельзя писать
	// то есть если buffeed канал, то спокойно дочитываем данные из буфера, например, через range
	// если unBuffered, то читаем значение, при последуюищх попытках поулчим 0, false
	fmt.Printf("got %d from unBuffered\n", <-unBuf)

}

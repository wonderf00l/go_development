package main

import "fmt"

func greet(roc <-chan string) {
	fmt.Println("Hello " + <-roc + "!")
	// roc <- "awewa" не можем писать в канал внутри именно этой функции
}

func main() {
	fmt.Println("main() started")
	c := make(chan string)

	go greet(c)

	c <- "John"
	fmt.Println("main() stopped")
	{
		fmt.Println("main() started")
		c := make(chan string)

		// launch anonymous goroutine
		go func(c chan string) {
			fmt.Println("Hello " + <-c + "!")
		}(c)

		c <- "John"
		fmt.Println("main() stopped")
	}
}

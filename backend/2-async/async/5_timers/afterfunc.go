package main

import (
	"fmt"
	"time"
)

func sayHello() {
	fmt.Println("Hello World")
}

func main() {
	//timer := time.AfterFunc(1*time.Second, sayHello) вызов функции SayHello через 1 сек
	//
	//fmt.Scanln()
	//timer.Stop()

	timer := time.NewTimer(2 * time.Second)
	t := <-timer.C // почему не возник deadlock: C - буферизированный канал(chan *time.Time, 1)

	fmt.Println("Timer", t)

	t = <-time.After(1 * time.Second) // сразу возвращает канал
	// time.After нельзя остановить
	fmt.Println("Time after", t)

}

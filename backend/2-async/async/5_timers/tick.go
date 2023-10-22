package main

import (
	"fmt"
	"time"
)

func main() {
	ticker := time.NewTicker(time.Second)
	i := 0
	for tickTime := range ticker.C { // (chan Time, 1)
		i++
		// fmt.Println(runtime.NumGoroutine())
		fmt.Println("step", i, "time", tickTime)
		if i >= 5 {
			// надо останавливать, иначе потечет(то есть будет вечно тикать)
			ticker.Stop() // but doesn't close the channel, so need break
			break
		}
	}
	fmt.Println("total", i)
	// return

	// не может быть остановлен и собран сборщиком мусора
	// используйте если должен работать вечено
	c := time.Tick(time.Second)
	i = 0
	for tickTime := range c {
		i++
		fmt.Println("step", i, "time", tickTime)
		if i >= 5 {
			break
		}
	}

}

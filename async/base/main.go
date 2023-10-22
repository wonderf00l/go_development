package main

import (
	"fmt"
	"sync"
)

func main() {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		// wg.Add(1) -- если добавлять гоуртину внутри самой горутины, она может не успеть добавить себя в WaitGroup
		defer wg.Done()
		fmt.Println("hello from goroutine")
	}()
	fmt.Println("hello form Main goroutine") // без wg и sleep почти всегда выводится только 'hello form main goroutine'
	wg.Wait()
	// time.Sleep(1 * time.Second) -- со слипом скорее всего успеем дождаться горутины
}

package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	runtime.GOMAXPROCS(1)
	var counters = map[int]int{} // мапы не потокобезопасны, для безопасного обновления таких структур используем mutex
	mu = &sync.Mutex{}           // гарантия что mu будет только 1 и не будет копироваться
	for i := 0; i < 5; i++ {
		go func(counters map[int]int, th int) {
			for j := 0; j < 5; j++ {
				mu.Lock()
				counters[th*10+j]++
				mu.Unlock()
			}
		}(counters, i)
	}
	fmt.Scanln()
	mu.Lock()
	fmt.Println("counters result", counters) // конкурентное чтение
	mu.Unlock()
}

// -race --> race detector

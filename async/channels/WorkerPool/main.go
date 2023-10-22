package main

import (
	"fmt"
	"time"
)

// worker than make squares
func sqrWorker(tasks <-chan int, results chan<- int, id int) {
	for num := range tasks {
		time.Sleep(time.Millisecond) // simulating blocking task
		fmt.Printf("[worker %v] Sending result by worker %v\n", id, id)
		results <- num * num
	}
}

func main() {
	fmt.Println("[main] main() started")

	tasks := make(chan int, 10)
	results := make(chan int, 10)

	// launching 3 worker goroutines
	for i := 0; i < 3; i++ {
		go sqrWorker(tasks, results, i)
	}

	// passing 5 tasks
	for i := 0; i < 5; i++ {
		tasks <- i * 2 // non-blocking as buffer capacity is 10
	}

	fmt.Println("[main] Wrote 5 tasks")

	// closing tasks
	close(tasks)

	// receving results from all workers
	for i := 0; i < 5; i++ {
		result := <-results // blocking because buffer is empty
		fmt.Println("[main] Result", i, ":", result)
	}

	fmt.Println("[main] main() stopped")
}

/*
Приведенный пример достаточно большой, но прекрасно объясняет,
как несколько горутин могут извлекать данные из канала и
выполнять свою работу. Горутины весьма эффективны,
когда они могут блокироваться.
Если убрать вызов time.Sleep(),
то только одна горутина будет выполняться,
так как другие горутины не будут запланированы,
до тех пор пока цикл не закончится и горутина не завершится.
*/

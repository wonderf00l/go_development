package main

import (
	"fmt"
	"time"
)

func main() {
	// ch := make(chan int)
	// ch <- 2     // пишем в канал, котоорый никто не читает параллельно -> Deadlock
	//num := <-ch // пытаемся считать из канала, в который никто не пишет, -> deadlock
	// fmt.Println(num)

	tick := time.Tick(100 * time.Millisecond)
	boom := time.After(500 * time.Millisecond)
	fmt.Println(len(tick), cap(tick))
	for {
		select {
		case <-tick:
			fmt.Println("tick.")
		case <-boom:
			fmt.Println("BOOM!")
			return
		default:
			fmt.Println("    .")
			time.Sleep(50 * time.Millisecond)
		}
	}

}

// package main

// import (
// 	"fmt"
// )

// func squares(c chan int) {
// 	for i := 0; i <= 3; i++ {
// 		num := <-c
// 		fmt.Println(num * num)
// 	}
// }

// func main() {
// 	fmt.Println("main() started")
// 	c := make(chan int, 3)

// 	go squares(c)
// 	// time.Sleep(1 * time.Second) min ~1 сек нужна горутине, чтобы успеть считать данные из канала
// 	c <- 1
// 	c <- 2
// 	c <- 3

// 	fmt.Println("main() stopped")
// }

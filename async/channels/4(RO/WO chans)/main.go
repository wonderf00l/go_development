package main

import (
	"fmt"
	"runtime"
)

func square(c chan int) {
	fmt.Println("[square] reading")
	num := <-c
	fmt.Printf("square got: %d\n", num)
	c <- num * num
	fmt.Printf("square write value to chan: %d\n", num)
}

func cube(c chan int) {
	fmt.Println("[cube] reading")
	num := <-c
	fmt.Printf("cube got: %d\n", num)
	c <- num * num * num
	fmt.Printf("cube write value to chan: %d\n", num)
}

func main() {
	fmt.Println("[main] main() started")

	squareChan := make(chan int)
	cubeChan := make(chan int)

	go square(squareChan)
	go cube(cubeChan) // шедулер инициализирует горутины и закидывает их в FIFO структуру, поэтому cube() начнет работать первой
	// time.Sleep(1 * time.Second)         // Sleep pauses the current goroutine for at least the duration d. A negative or zero duration causes Sleep to return immediately. то есть горутина блокируется, в этот момент инвокаются другие горутины и выполняют свою работу(в даннном случае выводят принты)
	fmt.Println("existing goroutines num: ", runtime.NumGoroutine()) // 3
	testNum := 3
	fmt.Println("[main] sent testNum to squareChan")

	squareChan <- testNum // здесь будет блокировка, шедулер инвокнет сначала cube
	// она выведет начальные принты, но заблочится на чтении из канала
	// далее шедулер попробует выполнить square,она отпринтит, прочитает из канала
	// и заблочится уже на этапе записи в канал

	fmt.Println("[main] resuming")
	fmt.Println("[main] sent testNum to cubeChan")

	fmt.Println("existing goroutines num: ", runtime.NumGoroutine()) // 3

	cubeChan <- testNum // снова блокировка, вызываем незаблокированные горуитины,
	//сейчас это толкьо cube,  принтов уже не будет, блок на записи в канал

	fmt.Println("[main] resuming")
	fmt.Println("[main] reading from channels")

	squareVal, cubeVal := <-squareChan, <-cubeChan                   // блокировка, пробуждаем каждую из горутин
	sum := squareVal + cubeVal                                       // тут они уже завершили работу, main дорабатывает свое и завершается
	fmt.Println("existing goroutines num: ", runtime.NumGoroutine()) // 1

	fmt.Println("[main] sum of square and cube of", testNum, " is", sum)
	fmt.Println("[main] main() stopped")
}

// итог: порядок выполнения горутин не определен
// полагатсья на него не стоит
// когда-то выводится сразу square reading и cuve reading
// когда то cube reading похже и тд
// главное, чтобы не было взаимных блокировок

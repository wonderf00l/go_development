package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

/*
AddT
LoadT
StoreT
SwapT
CompareAndSwapT
atomic.Value
*/

/*
T:
int32
int64
uint32
uint64
uintptr
*/

func atomicValue() {
	var value atomic.Value

	value.Store(1)
	value.Load()
	value.Swap(2)
	value.CompareAndSwap(1, 2) // false
}

func main() {
	var intVar int32 = 0

	fmt.Println(atomic.AddInt32(&intVar, 1)) // intVAR == 1
	fmt.Println(atomic.LoadInt32(&intVar))   // 1

	var newVal int32 = 2

	atomic.StoreInt32(&intVar, newVal)        // intVAR == 2
	fmt.Println(atomic.SwapInt32(&intVar, 3)) // got: 2, intVar == 3

	wg := &sync.WaitGroup{}

	wg.Add(100)
	for i := 0; i != 100; i++ {
		go func(id int) {
			defer wg.Done()
			if !atomic.CompareAndSwapInt32(&intVar, 3, 4) { //  true, если получилось поменять 4 на 3
				return
			}
			fmt.Println("горутина, певрая успевшая заинкрементить: ", id)
		}(i)
	}
	// runtime.Gosched()
	wg.Wait()
}

// atomic для чисел - лучше, быстрее, оснвоа пакеты sync
// также позволяет рабоать с указателями

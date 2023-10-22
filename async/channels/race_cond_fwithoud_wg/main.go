package main

import (
	"fmt"
	"time"
)

var i int

func main() {
	for j := 0; j != 1000; j++ {
		go func(id int) {
			fmt.Println("inside worker ", id)
			i = i + 1
		}(j)
	}
	time.Sleep(1 * time.Second)
	fmt.Println("res:", i)
}

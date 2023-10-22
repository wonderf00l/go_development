package main

import (
	"fmt"

	prsonModule "github.com/go-park-mail-ru/lectures/1-basics/5_visibility/person"
)

func init() {
	fmt.Println("I am init in main")
}

func main() {
	person := 1
	fmt.Println("I am main", person, prsonModule.Person{})
}

// папка - название пакета (должно совпадать, эквивалентные сущности)
// с большой буквы - экспортируемое
// с мал буквы - только в рамках ВСЕГО пакета

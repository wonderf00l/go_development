package main

import (
	"fmt"
)

func foo(lol *[]string) {
	*lol = append(*lol, "exist")
}

func boo(kek []string) {
	kek = append(kek, "destroy")                    // значение пишется по указателю после уже существующих элементов
	fmt.Println("before return from function", kek) // меняется len, cap у копии --> при выходе len и cap остаются прежними --> неважно, что было записано после исходного участка памяти аргумента, к этой области памяти все равно нет доступа
} // а потому append для локальной копии не влияет на оригинал

func zoo(shrek []string) {
	shrek[0] = "changed" // под капотом обращаемся к значению через разыменование указателя --> значение изменится
}

func main() {
	test := []string{}
	fmt.Printf("0 create empty slice %v\n\n", test)
	test = append(test, "1", "2", "3", "4", "5")
	test = append(test, "6")

	fmt.Println("1 pass slice by value and append")
	boo(test)
	fmt.Println("after return from function", test)
	fmt.Println(test[6]) // не можем обратиться по id > len, то есть в неинициализированную ячейку, даже в рамках допустимого cap
	fmt.Println(len(test), cap(test))

	// fmt.Printf("\n2 pass slice by ref and append\n")
	// foo(&test)
	// fmt.Println("after return from function", test)

	// fmt.Printf("\n3 pass slice by value and change element\n")
	// zoo(test)
	// fmt.Println("after return from function", test)
}

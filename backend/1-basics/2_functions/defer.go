package main

import "fmt"

// defer - для recover, освобожжения ресурсов

func getSomeVars() string {
	fmt.Println("getSomeVars execution")
	return "getSomeVars result"
}

func withArgs(arg string) {
	fmt.Println("with_args")
}

/*
func main() {
	defer func() {
		fmt.Println("enty defer")
	}() // выполнится в конце

	for i := 0; i != 5; i++ {
		fmt.Println("not deferred print inside for") // выведется первым - i раз
		defer fmt.Println("deferred print inside for") // накопится i вызовов, вызовутся после выхода из блока for
	}
}

not deferred print inside for
not deferred print inside for
not deferred print inside for
not deferred print inside for
not deferred print inside for
deferred print inside for
deferred print inside for
deferred print inside for
deferred print inside for
deferred print inside for
enty defer

*/

func main() {
	defer fmt.Println("After work")

	// defer withArgs(getSomeVars()) -- getSomeVars() выполняется сразу при заходе в тело defer-функций

	defer func() { // через анонимную функцию - getsomevars вызовется в конце
		fmt.Println(getSomeVars()) // т.к. вызов getsomevars уже не является аргументом defer функции, который ИСПОЛНЯЕТСЯ СРАЗУ ПРИ ЗАХОДЕ В ФУНКЦИЮ
	}()

	// defer - в рамках конкретной функции/тела

	// for i := range(make([]int{})) {
	// 	// defer func ... - выполняется отложенно в РАМКАХ for
	// }

	fmt.Println("Some userful work")
}

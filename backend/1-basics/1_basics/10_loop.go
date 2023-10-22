package main

import "fmt"

func main() {
	// цикл без условия, while(true) OR for(;;;)
	for {
		fmt.Println("loop iteration")
		break
	}

	// цикл без условия, while(isRun)
	isRun := true
	for isRun {
		fmt.Println("loop iteration with condition")
		isRun = false
	}

	// цикл с условие и блоком инициализации
	for i := 0; i < 2; i++ {
		fmt.Println("loop iteration", i)
		if i == 1 {
			continue
		}
	}

	// операции по slice
	sl := []int{1, 2, 3}
	idx := 0

	for idx < len(sl) {
		fmt.Println("while-stype loop, idx:", idx, "value:", sl[idx])
		idx++
	}

	for i := 0; i < len(sl); i++ {
		fmt.Println("c-style loop", i, sl[i])
	}
	for idx := range sl { // при изменении слайса во время цикла range остается изначальный
		fmt.Println("range slice by index", sl[idx])
	}
	for idx, val := range sl {
		fmt.Println("range slice by idx-value", idx, val)
	}

	// операции по map(порядок значений не определен, будут выводиться по-разному)
	profile := map[int]string{1: "Vasily", 2: "Romanov"}

	for key := range profile {
		fmt.Println("range map by key", key)
	}

	for key, val := range profile {
		fmt.Println("range map by key-val", key, val)
	}

	for _, val := range profile {
		fmt.Println("range map by val", val)
	}

	str := "Привет, Мир! 💩" // обход не по байтам, а по символам
	for pos, char := range str {
		fmt.Printf("%#U at pos %d\n", char, pos)
	}

}

package main

import "fmt"

type Person struct {
	id      int
	Name    string
	Address string
}

type Account struct {
	Id      int
	Name    string
	Cleaner func(string) string
	Person  // встраивание 1-ой структуры в другую (в случае интерфейсов - композиция интерфейсов)
}

func main() {
	// полное объявление структуры
	var acc Account = Account{
		Id:   1,
		Name: "rvasily",
		Person: Person{
			id:      2,
			Name:    "Василий",
			Address: "Москва",
		},
	}
	fmt.Printf("%#v\n", acc)

	// короткое объявление структуры
	//acc.Person = Person{2, "Romanov Vasily", "Moscow"}

	// _ = &Person{
	// 	Name: "vas", // rest - default
	//}

	fmt.Printf("%#v\n", acc)

	fmt.Println(acc.Name)
	//fmt.Println(acc.Person.Name)
}

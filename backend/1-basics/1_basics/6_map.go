package main

import "fmt"

func main() {
	// инициализация при создании
	user := map[string]string{}

	user["test"] = "r"

	_ = map[string]byte{"one": 1} // или без map
	var _ map[string]byte
	// var m map[string]byte{"one": 1, "two": 2} ошибка: инициализация при создании через синтаксис выше

	// сразу с нужной ёмкостью
	profile := make(map[string]string, 10)

	// количество элементов
	mapLength := len(user)

	fmt.Printf("%d %+v\n", mapLength, profile)

	// если ключа нет - вернёт значение по умолчанию для типа
	mName := user["middleName"]
	fmt.Println("mName:", mName)

	// проверка на существование ключа
	mName, mNameExist := user["middleName"]
	fmt.Println("mName:", mName, "mNameExist:", mNameExist)

	// пустая переменная - только проверяем что ключ есть
	_, mNameExist2 := user["middleName"]
	fmt.Println("mNameExist2", mNameExist2)

	// удаление ключа
	delete(user, "lastName")
	fmt.Printf("%#v\n", user)

	//!!!
	// внутри мапы - hashmap, порядок элементов не определен
	// при добавлении может произойти рехеширование
}

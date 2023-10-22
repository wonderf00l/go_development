package main

import (
	"encoding/json"
	"fmt"
)

var jsonStr = `[
	{"id": 17, "username": "iivan", "phone": 0},
	{"id": "17", "address": "none", "company": "Mail.ru"},
	5,
	[]
]`

func main() {
	data := []byte(jsonStr)

	var users interface{}               // десериализуем в пустой интерфейс
	err := json.Unmarshal(data, &users) // в качестве 2-ого аргмента передаем указатель на объект, иначе создатся локальная копия на стеке, заполнится данными локально
	if err != nil {
		fmt.Println(err) // структура будет пустой, получем ошибку, если json поломанный
	}
	fmt.Printf("unpacked in empty interface:\n%#v\n\n", users)
	id := users.([]interface{})[0].(map[string]interface{})["id"].(float64)
	var i float64
	i = id - 2
	fmt.Printf("float value from data:\n", id+1)

	user2 := map[string]interface{}{
		"id":       42,
		"username": "rvasily",
	}
	var user2i interface{} = user2    // сначала создаем пустой интерфейс
	result, _ := json.Marshal(user2i) // затем упаковываем в json
	fmt.Printf("json string from map:\n %s\n", string(result))
}

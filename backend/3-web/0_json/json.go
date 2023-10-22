package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
)

type User struct {
	ID       int
	Username string
	phone    string // у Marshal метода нет доступа к приватным полям --> phone останется пустым после unmarshal и после marshal'а структуры в json объект/json строку
}

var jsonStr = `{"id": 42, "username": "rvasily", "phone": "123"}`

func main() {
	data := []byte(jsonStr)

	u := &User{}            // Unmarshall - десериализация json-данных(строка и тп) в какую-то встроенную сущность
	json.Unmarshal(data, u) // после инициализации структуры json-строкой названия полей - структуры, не json-строки
	fmt.Printf("struct:\n\t%#v\n\n", u)
	/*Unmarshal parses the JSON-encoded data
	and stores the result in the value pointed to by v.
	If v is nil or not a pointer, Unmarshal returns an
	InvalidUnmarshalError.*/

	u.phone = "987654321"
	result, err := json.Marshal(u) // сериализация(структуры в json-сущность)
	if err != nil {
		panic(err)
	}
	fmt.Printf("json string:\n\t%s\n", string(result))

	const jsonStream = `
	{"Name": "Ed", "Text": "Knock knock."}
	{"Name": "Sam", "Text": "Who's there?"}
	{"Name": "Ed", "Text": "Go fmt."}
	{"Name": "Sam", "Text": "Go fmt who?"}
	{"Name": "Ed", "Text": "Go fmt yourself!"}
`
	type Message struct {
		Name, Text string
	}
	dec := json.NewDecoder(strings.NewReader(jsonStream))
	for {
		var m Message
		if err := dec.Decode(&m); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s: %s\n", m.Name, m.Text)
	}

	// Encode - оборачиваем в Encoder каой-то Writer(через NewEncoder()),
	// далее будем писать в этот райтер JSON encoding данные(JSON-строку, структуры и тп,
	// все это под капотом маршалится(через marshal() - то есть конвертится из встроенных сущностей в json объекты и сразу записыватся в райтер)), в конце \n
	// Marshal же сериализует какую-то встроенную сущность в []byte - по факту JSON-encoded data
}

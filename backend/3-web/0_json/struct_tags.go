package main

import (
	"encoding/json"
	"fmt"
)

type User struct {
	ID       int    `json:"user_id,string"` // название поля в json-сущности, тип в josn-сущности
	Username string // название как в структуре, значение - тоже(+тип)
	Address  string `json:",omitempty"` // если поле будет пустым(с zero value), мы бы его не сериализовали
	Comnpany string `json:"-"`          // игнорируем поле при сериализации/десериализации
}

func main() {
	u := &User{
		ID:       42,
		Username: "rvasily",
		Address:  "test",
		Comnpany: "Mail.Ru Group",
	}
	result, _ := json.Marshal(u)
	fmt.Printf("json string: %s\n", string(result))
}

package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) { // при указании нескольких одинаковых параметров запроса (query/headers etc) их зачения агрегируются в слайс
	for key, value := range r.URL.Query() { // Query() - парсинг и итерация по параметрам запроса(query params)
		fmt.Println(key, value)
	}
}

func main() {
	http.HandleFunc("/", handler)

	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", nil)
}

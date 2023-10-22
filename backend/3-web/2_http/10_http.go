package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	// ResponseWriter - implements io.Writer
	w.Header().Set("Content-Type", "application/json") // сначала заголовки, потом тело
	fmt.Fprintln(w, "<h1>Привет, мир!</h1>")
	w.Write([]byte("{}"))
}

func main() {
	http.HandleFunc("/", handler)

	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", nil) // при указании nil используется DefaultServeMux, у которого есть мапа с эндпоинтами и хендлерами
}

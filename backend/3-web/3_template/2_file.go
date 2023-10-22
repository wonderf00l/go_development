package main

import (
	"fmt"
	"html/template" // уже html template
	"net/http"
)

type User struct {
	ID     int
	Name   string
	Active bool
}

func main() {
	tmpl := template.Must(template.ParseFiles("users.html")) // Must - если криво распарсили, то запаникуем

	users := []User{
		User{1, "Vasily", true},
		User{2, "<i>Ivan</i>", false}, // escape(экранирование) тегов
		User{3, "Dmitry", true},
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w,
			struct { // нужно именно struct {users []users}
				Users []User
			}{
				users,
			})
	})

	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", nil)
}

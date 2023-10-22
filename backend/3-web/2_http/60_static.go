package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`
		Hello World! <br />
		<img src="/data/img/gopher.png" />
	`)) // юзер пришел с таким Path до картинки
	// БЛАГОДАРЯ STRIPpREFIX до сервера дойдет только img/gopher.png, которое прибавится к ./static
}

/*<iframe />
<object />
<img />
<picture />
<embed />
<object />
<link />
<script />
<audio />
<video />
<track />
this tags can fetch external resources from the server
*/

func main() {
	http.HandleFunc("/", handler)

	staticHandler := http.StripPrefix(
		"/data/",
		http.FileServer(http.Dir("./static")),
	)
	// http.File
	http.Handle("/data/", staticHandler)

	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", nil)
}

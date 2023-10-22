package main

import (
	"fmt"
	"net/http"
)

func runServer(addr string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Addr:", addr, "URL:", r.URL.String())
		})

	mux.HandleFunc("/new", func(w http.ResponseWriter, r *http.Request) {
		serverNew := http.Server{
			Addr:    ":8082",
			Handler: mux,
		}
		fmt.Println("Starting New server at", ":8082")
		serverNew.ListenAndServe()
	}) // запуск нового сервера по эндпоинту /new, причем в него передается контекст того сервера(URL и тп), с которого пошли на эндпоинт

	server := http.Server{
		Addr:    addr,
		Handler: mux,
	}

	fmt.Println("starting server at", addr)
	server.ListenAndServe()
}

func main() {
	go runServer(":8081") // один сервер запускаем в горутине
	runServer(":8080")    // другой - в основном потоке, чтобы смогли заблокироваться, иначе горутины запустятся и сразу завершатся с главной горутиной
}

package main

import (
	"fmt"
	"net/http"

	_ "github.com/go-park-mail-ru/lectures/4-api/4_swagger/docs"
	handler "github.com/go-park-mail-ru/lectures/4-api/4_swagger/handlers"

	httpSwagger "github.com/swaggo/http-swagger"
)

// swag init

type myError struct {
	Status int
	Error  string
}

// @title Sample Project API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

func main() {

	fmt.Println("start server")
	http.HandleFunc("/docs/", httpSwagger.WrapHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("url: " + r.URL.String()))
	})
	http.HandleFunc("/users", handler.handleUsers)

	http.ListenAndServe(":8080", nil)

}

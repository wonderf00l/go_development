package main

import (
	"fmt"
	"net/http"
	"time"
)

var (
	client = http.Client{Timeout: time.Duration(time.Millisecond)}
)

type HTTPError struct {
	Code int
}

var e error = &ResourceError{}

type ResourceError struct {
	URL  string
	Err  error
	Code int
}

func (re *ResourceError) Error() string {
	return fmt.Sprintf(
		"Resource error: URL: %s, err: %v",
		re.URL,
		re.Err,
	)
}

func getRemoteResource() error {
	url := "http://127.0.0.1:9999/pages?id=123"
	_, err := client.Get(url)
	if err != nil {
		return &ResourceError{URL: url, Err: err}
	}
	return nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	err := getRemoteResource()
	if err != nil { // кастомные ошибки-структуры - не выражения, поэтому для них используется type assertion
		switch err.(type) { // type assertion в случае кастомных структур, реализующих error
		case *ResourceError:
			err := err.(*ResourceError)
			fmt.Printf("resource %s err: %s\n", err.URL, err.Err)
			http.Error(w, "remote resource error", 500)
		default:
			fmt.Printf("internal error: %+v\n", err)
			http.Error(w, "internal error", 500)
		}
		return
	}
	w.Write([]byte("all is OK"))
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", nil)
}

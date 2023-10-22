// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package gqlgen

type Author struct {
	Name string `json:"name"`
}

type Book struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Price  float64 `json:"price"`
	Author *Author `json:"author"`
}

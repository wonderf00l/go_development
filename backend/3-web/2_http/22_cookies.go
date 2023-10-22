package main

import (
	"fmt"
	"net/http"
	"time"
)

func mainPage(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session_id") // под капотом берем куки из хедера, парсим, вытаксииваем конкретный объект куки по ключу
	loggedIn := (err != http.ErrNoCookie)  // куки хранятся в браузере, после рестарта сервера все равно останемся в системе

	if loggedIn {
		fmt.Fprintln(w, `<a href="/logout">logout</a>`)
		fmt.Fprintln(w, "Welcome, "+session.Value)
	} else {
		fmt.Fprintln(w, `<a href="/login">login</a>`)
		fmt.Fprintln(w, "You need to login")
	}
}

func loginPage(w http.ResponseWriter, r *http.Request) { // сетим куку в expired, редиректим на mainpage
	expiration := time.Now().Add(10 * time.Hour) // выставляем location редиректа, заголовок Set-Cookie
	cookie := http.Cookie{
		Name:     "session_id",
		Value:    "Dmitry",
		Expires:  expiration,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", http.StatusFound)
}

func logoutPage(w http.ResponseWriter, r *http.Request) { // ставим куку, перенаправляем на основную страницу
	session, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, session)

	http.Redirect(w, r, "/", http.StatusFound)
}

func main() {
	http.HandleFunc("/login", loginPage)
	http.HandleFunc("/logout", logoutPage)
	http.HandleFunc("/", mainPage)

	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", nil)
}

package main

import (
	"fmt"
	"net/http"
	"time"
)

func mainPage(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session_id")
	// учебный пример! это не проверка авторизации!
	loggedIn := (err != http.ErrNoCookie)

	if loggedIn {
		fmt.Fprintln(w, `<a href="/logout">logout</a>`)
		fmt.Fprintln(w, "Welcome, "+session.Value)
	} else {
		fmt.Fprintln(w, `<a href="/login">login</a>`)
		fmt.Fprintln(w, "You need to login")
	}
}

func loginPage(w http.ResponseWriter, r *http.Request) {
	expiration := time.Now().Add(10 * time.Hour)
	cookie := http.Cookie{
		Name:    "session_id",
		Value:   "Dmitry",
		Expires: expiration,
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", http.StatusFound)
}

func logoutPage(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, session)
	http.Redirect(w, r, "/", http.StatusFound)
}

// -----------

func adminIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, `<a href="/">site index</a>`)
	fmt.Fprintln(w, "Admin main page")
}

func panicPage(w http.ResponseWriter, r *http.Request) {
	panic("this must me recovered")
}

// -----------

func pageWithAllChecks(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("recovered", err)
			http.Error(w, "Internal server error", 500)
		}
	}()
	defer func(start time.Time) {
		fmt.Printf("[%s] %s, %s %s\n",
			r.Method, r.RemoteAddr, r.URL.Path, time.Since(start))
	}(time.Now())

	_, err := r.Cookie("session_id")
	// учебный пример! это не проверка авторизации!
	if err != nil {
		fmt.Println("no auth at", r.URL.Path)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	// your logic
}

// -----------

func adminAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("adminAuthMiddleware", r.URL.Path) // привязываем мидлвару к урлу
		_, err := r.Cookie("session_id")
		// учебный пример! это не проверка авторизации!
		if err != nil {
			fmt.Println("no auth at", r.URL.Path)
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func accessLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("accessLogMiddleware", r.URL.Path)
		start := time.Now()
		next.ServeHTTP(w, r)
		fmt.Printf("[%s] %s, %s %s\n",
			r.Method, r.RemoteAddr, r.URL.Path, time.Since(start))
	})
}

func panicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("panicMiddleware", r.URL.Path)
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("recovered", err)
				http.Error(w, "Internal server error", 500)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// -----------

func main() {
	adminMux := http.NewServeMux() // отдельный мультиплексер для админки
	adminMux.HandleFunc("/admin/", adminIndex)
	adminMux.HandleFunc("/admin/panic", panicPage)
	// http.ServeMux - реализует http.Handler

	// set middleware
	adminHandler := adminAuthMiddleware(adminMux)

	siteMux := http.NewServeMux()
	siteMux.Handle("/admin/", adminHandler) // цепляем кластер хендлеров накрытых мидлварй к новому муксу, его будем юзать в ServeHTTP
	siteMux.HandleFunc("/login", loginPage)
	siteMux.HandleFunc("/logout", logoutPage)
	siteMux.HandleFunc("/", mainPage)

	// set middleware
	siteHandler := accessLogMiddleware(siteMux)
	siteHandler = panicMiddleware(siteHandler) // должна ловить паники во всех мидлварах, если будет до внедреении логов, паника не будет записана в логи
	// сервис не упадет после паники, отдадим 500

	// логирование -> проверка auth ->

	/*
		if path: /
		panicMiddleware /
		accessLogMiddleware /
		[GET] 127.0.0.1:46512, / 6.983µs
	*/

	/*
		if path: /admin/
		panicMiddleware /admin/
		accessLogMiddleware /admin/
		adminAuthMiddleware /admin/
		no auth at /admin/
		[GET] 127.0.0.1:35490, /admin/ 46.177µs
	*/

	/*
		если с авторизацией:
		path: /login
		panicMiddleware /login
		accessLogMiddleware /login
		[GET] 127.0.0.1:51228, /login 193.803µs
				|
				|
				V
		path: /admin/panic
		[GET] 127.0.0.1:39064, /admin/ 21.249µs
		panicMiddleware /admin/panic
		accessLogMiddleware /admin/panic
		adminAuthMiddleware /admin/panic
		recovered this must me recovered
	*/

	// если нет авторищации, то до panic() не дойдем, не пропустит auth middleware

	/*
		Слои:
		1. authMiddlware - проверка auth, дергаем исходный контроллер
		2. accessLogMidleware - инвокается слоем 3, пишутс логи, дергается ServeHTTP слоя 1
		3. panicMiddleware - сначала вызовется функция внешнего слоя, дернем ServeHTTP слоя 2, тут же отложенный recover()
	*/

	fmt.Println("starting server at http://127.0.0.1:8080")
	http.ListenAndServe(":8080", siteHandler)
}

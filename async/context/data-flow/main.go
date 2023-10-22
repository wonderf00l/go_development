package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
)

// передача данных через контекст - антипаттерн, должно хватит аргументов
// но в случае, если нужно прбросить какие-то данные во внутренний хендер запроса, можно использовать

// https://blog.ildarkarymov.ru/posts/context-guide/#%D0%BA%D0%BE%D0%B3%D0%B4%D0%B0-%D0%B8%D1%81%D0%BF%D0%BE%D0%BB%D1%8C%D0%B7%D0%BE%D0%B2%D0%B0%D1%82%D1%8C-%D0%BA%D0%BE%D0%BD%D1%82%D0%B5%D0%BA%D1%81%D1%82

type ctxKey string

const keyUserID ctxKey = "user_id"

func main() {
	mux := http.NewServeMux()

	mux.Handle("/restricted", authMiddleware(handleRestricted()))

	if err := http.ListenAndServe(":3000", mux); err != nil {
		log.Fatal(err)
	}
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Token")

		if token != "very-secret-token" {
			w.WriteHeader(http.StatusUnauthorized)
			io.WriteString(w, "invalid token")
			return
		}

		ctx := context.WithValue(r.Context(), keyUserID, 42)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func handleRestricted() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(keyUserID).(int)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, "internal error, try again later please")
			return
		}

		io.WriteString(w, fmt.Sprintf("hello, user #%d!", userID))
	})
}

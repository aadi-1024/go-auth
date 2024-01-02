package middlewarefunc

import (
	"log"
	"net/http"
)

func SessionAuth() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("auth-cookie")
			if err != nil {
				log.Println(err)
				http.Redirect(w, r, "/login", 303)
				return
			}
			if cookie.Value != "authenticated" {
				w.Write([]byte("couldnt login"))
				http.Redirect(w, r, "/login", 303)
			}

			next.ServeHTTP(w, r)
		})
	}
}

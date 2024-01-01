package middlewarefunc

import (
	"fmt"
	"net/http"
)

//func BasicAuthMiddleware(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		user, pass, ok := r.BasicAuth()
//		//Authorization header not present
//		if !ok {
//			basicAuthFail(w)
//			return
//		}
//		//if user != "user" || pass != "password" {
//		//	basicAuthFail(w)
//		//	return
//		//}
//		next.ServeHTTP(w, r)
//	})
//}

func BasicAuth(creds map[string]string, realm string) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, pass, ok := r.BasicAuth()
			//Auth header not present
			if !ok {
				basicAuthFail(w, realm)
				return
			}
			credp, ok := creds[user]
			if !ok || credp != pass {
				basicAuthFail(w, realm)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func basicAuthFail(w http.ResponseWriter, realm string) {
	w.Header().Add("WWW-Authenticate", fmt.Sprintf("Basic realm='%v'", realm))
	w.WriteHeader(http.StatusUnauthorized)
}

package main

import (
	"encoding/json"
	mdw "github.com/aadi-1024/go-auth/middleware"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"time"
)

func main() {
	mux := chi.NewMux()

	sessionAuthMux := chi.NewMux()
	sessionAuthMux.Use(mdw.SessionAuth())

	sessionAuthMux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		handle(w, r, "authenticated via session auth")
	})

	//setup basic auth
	basicAuthMux := chi.NewMux()
	basicAuthMux.Use(mdw.BasicAuth(map[string]string{
		"user":  "pass",
		"admin": "strongpass",
	}, "test-realm"))
	basicAuthMux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		handle(w, r, "Authenticated via basic auth")
	})

	mux.Mount("/basic-auth", basicAuthMux)
	mux.Mount("/session", sessionAuthMux)
	mux.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		if !r.Form.Has("passcode") {
			http.Error(w, "no password provided", 400)
		}
		if r.Form.Get("passcode") == "1234" {
			http.SetCookie(w, &http.Cookie{
				Name:    "auth-cookie",
				Value:   "authenticated",
				Expires: time.Now().Add(time.Hour),
			})
			w.WriteHeader(200)
		}
	})
	mux.Get("/login", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`
		<html>
			<head>
				<body>
					<form action="/login" method="post">
						<label for="passcode">Enter Passcode</label>
						<input type="password" name="passcode">
						<button type="submit">Submit</button>
					</form>
				</body>
			</head>
		</html>
		`))
	})
	if err := http.ListenAndServe(":8080", mux); err != nil {

		log.Fatalln(err)
	}
}

type jsonPayload struct {
	Data string `json:"data"`
}

func handle(w http.ResponseWriter, r *http.Request, data string) {
	payload := &jsonPayload{
		Data: data,
	}
	w.Header().Set("Content-Type", "application/json")
	buf, _ := json.Marshal(payload)
	w.Write(buf)
}

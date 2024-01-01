package main

import (
	"encoding/json"
	mdw "github.com/aadi-1024/go-auth/middleware"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func main() {
	mux := chi.NewMux()

	basicAuthMux := chi.NewMux()
	basicAuthMux.Use(mdw.BasicAuth(map[string]string{
		"user":  "pass",
		"admin": "strongpass",
	}, "test-realm"))
	basicAuthMux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		handle(w, r, "Authenticated via basic auth")
	})

	mux.Mount("/basic-auth", basicAuthMux)

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

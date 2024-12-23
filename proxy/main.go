package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Use(NewReverseProxy("hugo", "1313").ReverseProxy)

	r.Route("/api", func(r chi.Router) {
		r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello from API"))
		})

	})

	http.ListenAndServe(":8080", r)
}

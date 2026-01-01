package main

import (
	"net/http"
	"time"

	"github.com/Robert076/doclane/backend/handlers/auth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})
	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", auth.LoginHandler)
		r.Post("/register", auth.RegisterHandler)
	})
	http.ListenAndServe(":8080", r)
}

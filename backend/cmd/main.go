package main

import (
	"net/http"
	"time"

	"github.com/Robert076/doclane/backend/handlers/auth"
	auth_middleware "github.com/Robert076/doclane/backend/handlers/auth/middleware"
	document_handler "github.com/Robert076/doclane/backend/handlers/documents"
	user_handler "github.com/Robert076/doclane/backend/handlers/users"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
)

func main() {
	r := chi.NewRouter()
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	handler := corsHandler.Handler(r)

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	r.Route("/api/auth", func(r chi.Router) {
		r.Post("/login", auth.LoginHandler)
		r.Post("/register", auth.RegisterHandler)
		r.Post("/logout", auth.LogoutHandler)
	})

	r.Route("/api", func(r chi.Router) {
		r.Use(auth_middleware.AuthGuard)
		r.Use(auth_middleware.MustBeActive)
		r.Route("/users", func(r chi.Router) {
			r.Get("/", user_handler.GetUsersHandler)
			r.Get("/me", user_handler.GetCurrentUserHandler)
			r.Get("/my-clients", user_handler.GetClientsByProfessionalHandler)
		})

		r.Route("/document-requests", func(r chi.Router) {
			r.Post("/", document_handler.AddDocumentRequestHandler)

			r.Get("/professional/documents", document_handler.GetDocumentRequestsByProfessionalHandler)
			r.Get("/client/documents", document_handler.GetDocumentRequestsByClientHandler)

			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", document_handler.GetDocumentRequestByIDHandler)
				r.Put("/status", document_handler.UpdateDocumentRequestStatusHandler)

				r.Route("/files", func(r chi.Router) {
					r.Get("/", document_handler.GetFilesByRequestHandler)
					r.Post("/", document_handler.AddDocumentHandler)
				})
			})
		})
	})

	http.ListenAndServe(":8080", handler)
}

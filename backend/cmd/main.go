package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Robert076/doclane/backend/handlers/auth"
	auth_middleware "github.com/Robert076/doclane/backend/handlers/auth/middleware"
	document_handler "github.com/Robert076/doclane/backend/handlers/documents"
	invitation_handler "github.com/Robert076/doclane/backend/handlers/invitation"
	user_handler "github.com/Robert076/doclane/backend/handlers/users"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
)

func main() {
	r := chi.NewRouter()

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
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
		_, err := w.Write([]byte("Hello World!"))
		if err != nil {
			log.Fatal(err)
		}
	})

	r.Route("/api/auth", func(r chi.Router) {
		r.Post("/login", auth.LoginHandler)
		r.Post("/register/professional", auth.RegisterProfessionalHandler)
		r.Post("/register/client", auth.RegisterClientHandler)
		r.Post("/logout", auth.LogoutHandler)
	})

	r.Post("/api/invitations/validate", invitation_handler.ValidateInvitationCodeHandler)

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
				r.Patch("/", document_handler.PatchDocumentRequestHandler)
				r.Route("/files", func(r chi.Router) {
					r.Get("/", document_handler.GetFilesByRequestHandler)
					r.Post("/", document_handler.AddDocumentHandler)
					r.Get("/{fileId}/presign", document_handler.GetFilePresignedURLHandler)
				})
			})
		})

		r.Route("/invitations", func(r chi.Router) {
			r.Post("/generate", invitation_handler.GenerateInvitationCodeHandler)
			r.Get("/my-codes", invitation_handler.GetMyInvitationCodesHandler)
			r.Delete("/{id}", invitation_handler.DeleteInvitationCodeHandler)
		})
	})

	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}

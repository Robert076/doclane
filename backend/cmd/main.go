package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Robert076/doclane/backend/handlers/auth"
	auth_middleware "github.com/Robert076/doclane/backend/handlers/auth/middleware"
	document_handler "github.com/Robert076/doclane/backend/handlers/documents"
	invitation_handler "github.com/Robert076/doclane/backend/handlers/invitation"
	template_handler "github.com/Robert076/doclane/backend/handlers/templates"
	user_handler "github.com/Robert076/doclane/backend/handlers/users"
	"github.com/aws/aws-lambda-go/lambda"
	chiadapter "github.com/awslabs/aws-lambda-go-api-proxy/chi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
)

func buildRouter() (http.Handler, *chi.Mux) {
	// hello from lambda oidc should fail
	r := chi.NewRouter()

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

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

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Stone cold healthy"))
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
			r.Get("/{id}", user_handler.GetUserByIDHandler)
			r.Post("/notify/{id}", user_handler.NotifyUserHandler)
			r.Get("/me", user_handler.GetCurrentUserHandler)
			r.Get("/my-clients", user_handler.GetClientsByProfessionalHandler)
			r.Post("/deactivate/{id}", user_handler.DeactivateUserHandler)
		})

		r.Route("/document-requests", func(r chi.Router) {
			r.Post("/", document_handler.AddDocumentRequestHandler)
			r.Get("/professional/my-requests", document_handler.GetDocumentRequestsByProfessionalHandler)
			r.Get("/client/my-requests", document_handler.GetDocumentRequestsByClientHandler)
			r.Patch("/expected-documents/{id}/status", document_handler.PatchExpectedDocumentStatusHandler)
			r.Get("/expected-documents/{id}/presign-example", document_handler.GetExamplePresignedURLHandler)

			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", document_handler.GetDocumentRequestByIDHandler)
				r.Patch("/", document_handler.PatchDocumentRequestHandler)
				r.Post("/archive", document_handler.CloseDocumentRequestHandler)
				r.Post("/unarchive", document_handler.ReopenDocumentRequestHandler)

				r.Route("/files", func(r chi.Router) {
					r.Get("/", document_handler.GetFilesByRequestHandler)
					r.Post("/", document_handler.AddDocumentHandler)
					r.Get("/{fileId}/presign", document_handler.GetFilePresignedURLHandler)
				})
			})

		})

		r.Route("/templates", func(r chi.Router) {
			r.Post("/", template_handler.AddTemplateHandler)
			r.Get("/", template_handler.GetTemplatesByProfessionalHandler)
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", template_handler.GetTemplateByIDHandler)
				r.Post("/instantiate", template_handler.InstantiateTemplateHandler)
				r.Route("/expected-documents", func(r chi.Router) {
					r.Get("/", template_handler.GetExpectedDocumentTemplatesByTemplateIDHandler)
					r.Post("/", template_handler.AddExpectedDocumentTemplateHandler)
					r.Delete("/{expectedDocId}", template_handler.DeleteExpectedDocumentTemplateHandler)
					r.Get("/{expectedDocId}/presign-example", template_handler.PresignExampleHandler)
				})
				r.Post("/archive", template_handler.CloseTemplateHandler)
				r.Post("/unarchive", template_handler.ReopenTemplateHandler)
			})
		})

		r.Route("/invitations", func(r chi.Router) {
			r.Post("/generate", invitation_handler.GenerateInvitationCodeHandler)
			r.Get("/my-codes", invitation_handler.GetMyInvitationCodesHandler)
			r.Delete("/{id}", invitation_handler.DeleteInvitationCodeHandler)
		})
	})

	return corsHandler.Handler(r), r
}

func main() {
	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
		_, mux := buildRouter()
		adapter := chiadapter.NewV2(mux)
		lambda.Start(adapter.ProxyWithContextV2)
	} else {
		handler, _ := buildRouter()
		if err := http.ListenAndServe(":8080", handler); err != nil {
			log.Fatal(err)
		}
	}
}

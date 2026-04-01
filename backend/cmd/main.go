package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Robert076/doclane/backend/handlers/auth"
	auth_middleware "github.com/Robert076/doclane/backend/handlers/auth/middleware"
	comment_handler "github.com/Robert076/doclane/backend/handlers/comments"
	department_handler "github.com/Robert076/doclane/backend/handlers/departments"
	insertadmin_handler "github.com/Robert076/doclane/backend/handlers/insert-admin"
	invitation_handler "github.com/Robert076/doclane/backend/handlers/invitation"
	request_handler "github.com/Robert076/doclane/backend/handlers/requests"
	template_handler "github.com/Robert076/doclane/backend/handlers/templates"
	user_handler "github.com/Robert076/doclane/backend/handlers/users"
	"github.com/aws/aws-lambda-go/lambda"
	chiadapter "github.com/awslabs/aws-lambda-go-api-proxy/chi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
)

func buildRouter() (http.Handler, *chi.Mux) {
	// hello from lambda :d
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
		_, err := w.Write([]byte("Healthy"))
		if err != nil {
			log.Fatal(err)
		}
	})

	r.Route("/api/auth", func(r chi.Router) {
		r.Post("/login", auth.LoginHandler)
		r.Post("/register", auth.RegisterHandler)
		r.Post("/logout", auth.LogoutHandler)
		r.Post("/insert-admin", insertadmin_handler.InsertAdminHandler)
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
			r.Post("/deactivate/{id}", user_handler.DeactivateUserHandler)
		})

		r.Route("/requests", func(r chi.Router) {
			r.Post("/", request_handler.AddRequestHandler)
			r.Get("/assignee/{id}", request_handler.GetRequestsByAssigneeHandler)
			r.Get("/department/{id}", request_handler.GetRequestsByDepartmentHandler)
			r.Post("/forward/{id}", request_handler.ForwardRequestToDepartmentHandler)
			r.Patch("/expected-documents/{id}/status", request_handler.PatchExpectedDocumentStatusHandler)
			r.Get("/expected-documents/{id}/presign-example", request_handler.GetExamplePresignedURLHandler)

			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", request_handler.GetRequestByIDHandler)
				r.Patch("/", request_handler.PatchRequestHandler)
				r.Post("/archive", request_handler.CloseRequestHandler)
				r.Post("/unarchive", request_handler.ReopenRequestHandler)
				r.Route("/comments", func(r chi.Router) {
					r.Get("/", comment_handler.GetCommentsByRequest)
					r.Get("/{commentID}", comment_handler.GetCommentByID)
					r.Post("/", comment_handler.AddCommentHandler)
				})
				r.Route("/files", func(r chi.Router) {
					r.Get("/", request_handler.GetFilesByRequestHandler)
					r.Post("/", request_handler.AddDocumentHandler)
					r.Get("/{fileId}/presign", request_handler.GetFilePresignedURLHandler)
				})
			})
		})

		r.Route("/templates", func(r chi.Router) {
			r.Post("/", template_handler.AddRequestTemplateWithDocumentsHandler)
			r.Get("/", template_handler.GetRequestTemplatesHandler)
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", template_handler.GetRequestTemplateByIDHandler)
				r.Patch("/", template_handler.PatchRequestTemplateHandler)
				r.Delete("/", template_handler.DeleteRequestTemplateHandler)
				r.Route("/expected-documents", func(r chi.Router) {
					r.Get("/", template_handler.GetExpectedDocumentTemplatesByRequestTemplateIDHandler)
					r.Delete("/{expectedDocId}", template_handler.DeleteExpectedDocumentTemplateHandler)
					r.Get("/{expectedDocId}/presign-example", template_handler.PresignExampleHandler)
				})
				r.Post("/archive", template_handler.CloseRequestTemplateHandler)
				r.Post("/unarchive", template_handler.ReopenRequestTemplateHandler)
			})
		})

		r.Route("/departments", func(r chi.Router) {
			r.Get("/", department_handler.GetAllDepartmentsHandler)
			r.Post("/", department_handler.CreateDepartmentHandler)
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

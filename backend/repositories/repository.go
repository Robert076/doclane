package repositories

import (
	"context"

	"github.com/Robert076/doclane/backend/models"
)

type IUserRepository interface {
	GetUsers(ctx context.Context, limit *int, offset *int, orderBy *string, order *string) ([]models.User, error)
	GetUserByID(ctx context.Context, id int) (models.User, error)
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
	AddUser(ctx context.Context, user models.User) (int, error)
}

type IDocumentRepository interface {
	AddDocumentRequest(ctx context.Context, req models.DocumentRequest) (int, error)
	GetDocumentRequestByID(ctx context.Context, id int) (models.DocumentRequest, error)
	GetDocumentRequestsByProfessional(ctx context.Context, professionalID int) ([]models.DocumentRequest, error)
	GetDocumentRequestsByClient(ctx context.Context, clientID int) ([]models.DocumentRequest, error)
	UpdateDocumentRequestStatus(ctx context.Context, id int, status string) error

	AddDocumentFile(ctx context.Context, file models.DocumentFile) (int, error)
	GetFilesByRequest(ctx context.Context, requestID int) ([]models.DocumentFile, error)
}

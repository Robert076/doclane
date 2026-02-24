package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/Robert076/doclane/backend/models"
)

type IUserRepository interface {
	GetUsers(ctx context.Context, limit *int, offset *int, orderBy *string, order *string, search *string) ([]models.User, error)
	GetUserByID(ctx context.Context, id int) (models.User, error)
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
	GetUsersByProfessionalID(ctx context.Context, professionalID int, limit *int, offset *int) ([]models.User, error)
	AddUser(ctx context.Context, user models.User) (int, error)
	NotifyUser(ctx context.Context, userId int, time time.Time) error
	DeactivateUser(ctx context.Context, userId int) error
}

type IDocumentRepository interface {
	GetDocumentRequestByID(ctx context.Context, id int) (models.DocumentRequestDTORead, error)
	GetDocumentRequestsByProfessional(ctx context.Context, professionalID int, search *string) ([]models.DocumentRequestDTORead, error)
	GetDocumentRequestsByProfessionalWithExpectedDocs(ctx context.Context, professionalID int, search *string) ([]models.DocumentRequestDTORead, error)
	GetDocumentRequestsByClient(ctx context.Context, clientID int, search *string) ([]models.DocumentRequestDTORead, error)
	GetDocumentRequestsByClientWithExpectedDocs(ctx context.Context, clientID int, search *string) ([]models.DocumentRequestDTORead, error)
	AddDocumentRequest(ctx context.Context, req models.DocumentRequest) (int, error)
	AddDocumentRequestWithTx(ctx context.Context, req models.DocumentRequest, transaction *sql.Tx) (int, error)
	UpdateDocumentRequestTitle(ctx context.Context, id int, newTitle string) error
	CloseDocumentRequest(ctx context.Context, id int) error

	AddDocumentFile(ctx context.Context, file models.DocumentFile) (int, error)
	GetFilesByRequest(ctx context.Context, requestID int) ([]models.DocumentFileDTORead, error)
	GetFileByID(ctx context.Context, id int) (models.DocumentFile, error)
	GetFileByIDExtended(ctx context.Context, id int) (models.DocumentFileDTOExtended, error)
	SetFileUploaded(ctx context.Context, id int) error
}

type IInvitationCodeRepository interface {
	GetInvitationCodeByCode(ctx context.Context, code string) (models.InvitationCode, error)
	GetInvitationCodeByID(ctx context.Context, id int) (models.InvitationCode, error)
	GetInvitationCodesByProfessional(ctx context.Context, professionalID int) ([]models.InvitationCode, error)
	CreateInvitationCode(ctx context.Context, code string, professionalID int, expiresAt *time.Time) error
	InvalidateCode(ctx context.Context, id int) error
	ReactivateCode(ctx context.Context, code string) error
	DeleteCode(ctx context.Context, id int) error
}

type IExpectedDocumentRepository interface {
	GetExpectedDocumentsByRequest(ctx context.Context, requestId int) ([]models.ExpectedDocument, error)
	AddExpectedDocumentToRequest(ctx context.Context, requestId int, expectedDocument models.ExpectedDocument) (int, error)
	AddExpectedDocumentToRequestWithTx(ctx context.Context, tx *sql.Tx, ed models.ExpectedDocument) error
	MarkAsUploaded(ctx context.Context, expectedDocumentID int) error
	DeleteExpectedDocumentFromRequest(ctx context.Context, requestId int, expectedDocumentId int) error
}

type ITxManager interface {
	WithTx(ctx context.Context, fn func(tx *sql.Tx) error) error
}

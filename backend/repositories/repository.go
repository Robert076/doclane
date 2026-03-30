package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/Robert076/doclane/backend/models"
)

type IUserRepo interface {
	GetUsers(ctx context.Context, limit *int, offset *int, orderBy *string, order *string, search *string) ([]models.User, error)
	GetUserByID(ctx context.Context, id int) (models.User, error)
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
	GetUsersByProfessionalID(ctx context.Context, professionalID int, limit *int, offset *int) ([]models.User, error)
	AddUser(ctx context.Context, user models.User) (int, error)
	NotifyUser(ctx context.Context, userId int, time time.Time) error
	DeactivateUser(ctx context.Context, userId int) error
}

type IRequestRepo interface {
	GetRequestByID(ctx context.Context, id int) (models.RequestDTORead, error)
	GetRequestsByProfessional(ctx context.Context, professionalID int, search *string) ([]models.RequestDTORead, error)
	GetRequestsByProfessionalWithExpectedDocs(ctx context.Context, professionalID int, search *string) ([]models.RequestDTORead, error)
	GetRequestsByClient(ctx context.Context, clientID int, search *string) ([]models.RequestDTORead, error)
	GetRequestsByClientWithExpectedDocs(ctx context.Context, clientID int, search *string) ([]models.RequestDTORead, error)
	AddRequest(ctx context.Context, req models.Request) (int, error)
	AddRequestWithTx(ctx context.Context, req models.Request, transaction *sql.Tx) (int, error)
	UpdateRequestTitle(ctx context.Context, id int, newTitle string) error
	CloseRequest(ctx context.Context, id int) error
	ReopenRequest(ctx context.Context, id int) error

	AddDocument(ctx context.Context, file models.Document) (int, error)
	GetFilesByRequest(ctx context.Context, requestID int) ([]models.DocumentDTORead, error)
	GetFileByID(ctx context.Context, id int) (models.Document, error)
	GetFileByIDExtended(ctx context.Context, id int) (models.DocumentDTOExtended, error)
	SetFileUploaded(ctx context.Context, id int) error
}

type IRequestTemplateRepo interface {
	GetRequestTemplatesByProfessionalID(ctx context.Context, professionalID int) ([]models.RequestTemplate, error)
	GetRequestTemplateByID(ctx context.Context, id int) (models.RequestTemplate, error)
	AddRequestTemplate(ctx context.Context, tmp models.RequestTemplate) (int, error)
	AddRequestTemplateWithTx(ctx context.Context, tx *sql.Tx, tmp models.RequestTemplate) (int, error) // new
	PatchRequestTemplate(ctx context.Context, templateID int, tmp models.RequestTemplateDTOPatch) error
	CloseRequestTemplate(ctx context.Context, id int) error
	ReopenRequestTemplate(ctx context.Context, id int) error
	DeleteRequestTemplate(ctx context.Context, id int) error
}

type IExpectedDocumentTemplateRepo interface {
	GetByRequestTemplateID(ctx context.Context, templateID int) ([]models.ExpectedDocumentTemplate, error)
	GetByDocumentID(ctx context.Context, id int) (models.ExpectedDocumentTemplate, error)
	Add(ctx context.Context, t models.ExpectedDocumentTemplate) (int, error)
	AddWithTx(ctx context.Context, tx *sql.Tx, t models.ExpectedDocumentTemplate) (int, error) // new
	DeleteByID(ctx context.Context, id int) error
}

type IInvitationCodeRepo interface {
	GetInvitationCodeByCode(ctx context.Context, code string) (models.InvitationCode, error)
	GetInvitationCodeByID(ctx context.Context, id int) (models.InvitationCode, error)
	GetInvitationCodesByProfessional(ctx context.Context, professionalID int) ([]models.InvitationCode, error)
	CreateInvitationCode(ctx context.Context, code string, professionalID int, expiresAt *time.Time) error
	InvalidateCode(ctx context.Context, id int) error
	ReactivateCode(ctx context.Context, code string) error
	DeleteCode(ctx context.Context, id int) error
}

type IExpectedDocumentRepo interface {
	GetExpectedDocumentByID(ctx context.Context, id int) (models.ExpectedDocument, error)
	GetExpectedDocumentsByRequest(ctx context.Context, requestId int) ([]models.ExpectedDocument, error)
	AddExpectedDocumentToRequest(ctx context.Context, requestId int, expectedDocument models.ExpectedDocument) (int, error)
	AddExpectedDocumentToRequestWithTx(ctx context.Context, tx *sql.Tx, ed models.ExpectedDocument) error
	UpdateExpectedDocumentStatus(ctx context.Context, documentID int, status string, rejectionReason *string) error
	DeleteExpectedDocumentFromRequest(ctx context.Context, requestId int, expectedDocumentId int) error
}

type IRequestCommentRepo interface {
	GetCommentsByRequestID(ctx context.Context, requestID int) ([]models.RequestCommentDTO, error)
	GetCommentByID(ctx context.Context, commentID int) (models.RequestCommentDTO, error)
	AddComment(ctx context.Context, comment models.RequestComment) (int, error)
	GetLastCommentFromUser(ctx context.Context, userID int) (models.RequestComment, error)
}

type ITxManager interface {
	WithTx(ctx context.Context, fn func(tx *sql.Tx) error) error
}

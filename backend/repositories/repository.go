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
	GetUsersByDepartment(ctx context.Context, departmentID int) ([]models.User, error)
	AddUser(ctx context.Context, user models.User) (int, error)
	NotifyUser(ctx context.Context, userId int, time time.Time) error
	DeactivateUser(ctx context.Context, userId int) error
	UpdatePassword(ctx context.Context, userID int, hashedPassword string) error
	UpdateUserProfile(ctx context.Context, userID int, dto models.UserProfilePatchDTO) error
	UpdateUserDepartment(ctx context.Context, userID int, departmentID int) error
}

type IRequestRepo interface {
	GetAllRequests(ctx context.Context, search *string) ([]models.RequestDTORead, error)
	GetRequestByID(ctx context.Context, id int) (models.RequestDTORead, error)
	GetRequestsByAssigneeWithExpectedDocs(ctx context.Context, assignee int, search *string) ([]models.RequestDTORead, error)
	GetRequestsByDepartment(ctx context.Context, departmentID int, search *string) ([]models.RequestDTORead, error)
	GetRequestsByDepartmentWithExpectedDocs(ctx context.Context, departmentID int, search *string) ([]models.RequestDTORead, error)
	GetArchivedRequests(ctx context.Context, search *string) ([]models.RequestDTORead, error)
	GetArchivedRequestsByDepartment(ctx context.Context, departmentID int, search *string) ([]models.RequestDTORead, error)
	GetCancelledRequests(ctx context.Context, search *string) ([]models.RequestDTORead, error)
	GetCancelledRequestsByDepartment(ctx context.Context, departmentID int, search *string) ([]models.RequestDTORead, error)
	GetDueRecurringRequests(ctx context.Context) ([]models.RequestDTORead, error)
	UpdateNextDueAt(ctx context.Context, requestID int, nextDueAt time.Time) error
	ClaimRequest(ctx context.Context, requestID int, userID int) error
	UnclaimRequest(ctx context.Context, requestID int) error
	AddRequest(ctx context.Context, req models.Request) (int, error)
	AddRequestWithTx(ctx context.Context, req models.Request, tx *sql.Tx) (int, error)
	UpdateRequestTitle(ctx context.Context, id int, newTitle string) error
	CancelRequest(ctx context.Context, id int) error
	CloseRequest(ctx context.Context, id int) error
	ReopenRequest(ctx context.Context, id int) error
	AddDocument(ctx context.Context, file models.Document) (int, error)
	GetFilesByRequest(ctx context.Context, requestID int) ([]models.DocumentDTORead, error)
	GetFileByID(ctx context.Context, id int) (models.Document, error)
	GetFileByIDExtended(ctx context.Context, id int) (models.DocumentDTOExtended, error)
	SetFileUploaded(ctx context.Context, id int) error
}

type IRequestTemplateRepo interface {
	GetRequestTemplatesByDepartment(ctx context.Context, departmentID int) ([]models.RequestTemplateDTORead, error)
	GetAllRequestTemplates(ctx context.Context) ([]models.RequestTemplateDTORead, error)
	GetRequestTemplateByID(ctx context.Context, id int) (models.RequestTemplateDTORead, error)
	AddRequestTemplate(ctx context.Context, tmp models.RequestTemplate) (int, error)
	AddRequestTemplateWithTx(ctx context.Context, tx *sql.Tx, tmp models.RequestTemplate) (int, error)
	CloseRequestTemplate(ctx context.Context, id int) error
	ReopenRequestTemplate(ctx context.Context, id int) error
	DeleteRequestTemplate(ctx context.Context, id int) error
	PatchRequestTemplate(ctx context.Context, templateID int, dto models.RequestTemplateDTOPatch) error
}

type IExpectedDocumentTemplateRepo interface {
	GetByRequestTemplateID(ctx context.Context, templateID int) ([]models.ExpectedDocumentTemplate, error)
	GetByDocumentID(ctx context.Context, id int) (models.ExpectedDocumentTemplate, error)
	Add(ctx context.Context, t models.ExpectedDocumentTemplate) (int, error)
	AddWithTx(ctx context.Context, tx *sql.Tx, t models.ExpectedDocumentTemplate) (int, error) // new
	DeleteByID(ctx context.Context, id int) error
}

type IDepartmentRepo interface {
	GetAllDepartments(ctx context.Context) ([]models.Department, error)
	GetDepartmentByID(ctx context.Context, id int) (models.Department, error)
	CreateDepartment(ctx context.Context, name string) (int, error)
}

type IInvitationCodeRepo interface {
	GetInvitationCodesByDepartment(ctx context.Context, departmentID int) ([]models.InvitationCode, error)
	GetInvitationCodeByCode(ctx context.Context, code string) (models.InvitationCodeReadDTO, error)
	GetInvitationCodesByCreator(ctx context.Context, createdBy int) ([]models.InvitationCode, error)
	GetInvitationCodeByID(ctx context.Context, id int) (models.InvitationCode, error)
	CreateInvitationCode(ctx context.Context, departmentID int, code string, createdBy int, expiresAt *time.Time) error
	InvalidateCode(ctx context.Context, id int) error
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

type IStatsRepo interface {
	GetStats(ctx context.Context) (*models.Stats, error)
}

type ITxManager interface {
	WithTx(ctx context.Context, fn func(tx *sql.Tx) error) error
}

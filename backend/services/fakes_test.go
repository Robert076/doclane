package services

import (
	"context"
	"database/sql"
	"time"

	"github.com/Robert076/doclane/backend/models"
	"github.com/Robert076/doclane/backend/repositories"
)

type fakeRequestRepo struct {
	repositories.IRequestRepo

	getRequestByID func(ctx context.Context, id int) (models.RequestDTORead, error)
}

func (f *fakeRequestRepo) GetRequestByID(ctx context.Context, id int) (models.RequestDTORead, error) {
	return f.getRequestByID(ctx, id)
}

type fakeCommentRepo struct {
	repositories.IRequestCommentRepo

	getLastCommentFromUser func(ctx context.Context, userID int) (models.RequestComment, error)
	addComment             func(ctx context.Context, comment models.RequestComment) (int, error)

	addedComment *models.RequestComment
}

func (f *fakeCommentRepo) GetLastCommentFromUser(ctx context.Context, userID int) (models.RequestComment, error) {
	return f.getLastCommentFromUser(ctx, userID)
}

func (f *fakeCommentRepo) AddComment(ctx context.Context, comment models.RequestComment) (int, error) {
	f.addedComment = &comment
	return f.addComment(ctx, comment)
}

type fakeInvitationRepo struct {
	repositories.IInvitationCodeRepo

	getByCreator         func(ctx context.Context, createdBy int) ([]models.InvitationCode, error)
	getByCode            func(ctx context.Context, code string) (models.InvitationCodeReadDTO, error)
	getByID              func(ctx context.Context, id int) (models.InvitationCode, error)
	createInvitationCode func(ctx context.Context, departmentID int, code string, createdBy int, expiresAt *time.Time) error
	deleteCode           func(ctx context.Context, id int) error
	deletedIDs           []int
	createdCode          string
	createdExpiresAt     *time.Time
}

func (f *fakeInvitationRepo) GetInvitationCodesByCreator(ctx context.Context, createdBy int) ([]models.InvitationCode, error) {
	return f.getByCreator(ctx, createdBy)
}

func (f *fakeInvitationRepo) GetInvitationCodeByCode(ctx context.Context, code string) (models.InvitationCodeReadDTO, error) {
	return f.getByCode(ctx, code)
}

func (f *fakeInvitationRepo) GetInvitationCodeByID(ctx context.Context, id int) (models.InvitationCode, error) {
	return f.getByID(ctx, id)
}

func (f *fakeInvitationRepo) CreateInvitationCode(ctx context.Context, departmentID int, code string, createdBy int, expiresAt *time.Time) error {
	f.createdCode = code
	f.createdExpiresAt = expiresAt
	return f.createInvitationCode(ctx, departmentID, code, createdBy, expiresAt)
}

func (f *fakeInvitationRepo) DeleteCode(ctx context.Context, id int) error {
	f.deletedIDs = append(f.deletedIDs, id)
	return f.deleteCode(ctx, id)
}

type fakeDepartmentRepo struct {
	repositories.IDepartmentRepo

	getDepartmentByID func(ctx context.Context, id int) (models.DepartmentDTORead, error)
}

func (f *fakeDepartmentRepo) GetDepartmentByID(ctx context.Context, id int) (models.DepartmentDTORead, error) {
	return f.getDepartmentByID(ctx, id)
}

type fakeUserRepo struct {
	repositories.IUserRepo

	getUserByID func(ctx context.Context, id int) (models.User, error)
}

func (f *fakeUserRepo) GetUserByID(ctx context.Context, id int) (models.User, error) {
	return f.getUserByID(ctx, id)
}

type fakeTemplateRepo struct {
	repositories.IRequestTemplateRepo

	getRequestTemplateByID func(ctx context.Context, id int) (models.RequestTemplateDTORead, error)
}

func (f *fakeTemplateRepo) GetRequestTemplateByID(ctx context.Context, id int) (models.RequestTemplateDTORead, error) {
	return f.getRequestTemplateByID(ctx, id)
}

type fakeExpectedDocTmplRepo struct {
	repositories.IExpectedDocumentTemplateRepo

	getByRequestTemplateID func(ctx context.Context, templateID int) ([]models.ExpectedDocumentTemplate, error)
}

func (f *fakeExpectedDocTmplRepo) GetByRequestTemplateID(ctx context.Context, templateID int) ([]models.ExpectedDocumentTemplate, error) {
	return f.getByRequestTemplateID(ctx, templateID)
}

type fakeTxManager struct {
	fn func(ctx context.Context, fn func(tx *sql.Tx) error) error
}

func (f *fakeTxManager) WithTx(ctx context.Context, fn func(tx *sql.Tx) error) error {
	if f.fn != nil {
		return f.fn(ctx, fn)
	}
	return fn(nil)
}

type fakeRequestRepoFull struct {
	repositories.IRequestRepo

	getRequestByID   func(ctx context.Context, id int) (models.RequestDTORead, error)
	addRequestWithTx func(ctx context.Context, req models.Request, tx *sql.Tx) (int, error)
	claimRequest     func(ctx context.Context, requestID, userID int) error
	cancelRequest    func(ctx context.Context, id int) error

	claimedRequestID int
	claimedByUser    int
	cancelledID      int
}

func (f *fakeRequestRepoFull) GetRequestByID(ctx context.Context, id int) (models.RequestDTORead, error) {
	return f.getRequestByID(ctx, id)
}

func (f *fakeRequestRepoFull) AddRequestWithTx(ctx context.Context, req models.Request, tx *sql.Tx) (int, error) {
	return f.addRequestWithTx(ctx, req, tx)
}

func (f *fakeRequestRepoFull) ClaimRequest(ctx context.Context, requestID, userID int) error {
	f.claimedRequestID = requestID
	f.claimedByUser = userID
	if f.claimRequest != nil {
		return f.claimRequest(ctx, requestID, userID)
	}
	return nil
}

func (f *fakeRequestRepoFull) CancelRequest(ctx context.Context, id int) error {
	f.cancelledID = id
	if f.cancelRequest != nil {
		return f.cancelRequest(ctx, id)
	}
	return nil
}

type fakeExpectedDocRepo struct {
	repositories.IExpectedDocumentRepo

	added []models.ExpectedDocument
}

func (f *fakeExpectedDocRepo) AddExpectedDocumentToRequestWithTx(ctx context.Context, tx *sql.Tx, ed models.ExpectedDocument) error {
	f.added = append(f.added, ed)
	return nil
}

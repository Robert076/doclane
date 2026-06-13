package services

import (
	"context"
	"time"

	"github.com/Robert076/doclane/backend/models"
	"github.com/Robert076/doclane/backend/repositories"
)

// The fakes below embed the repository interfaces so they satisfy the full
// interface contract without implementing every method. Each test overrides
// only the behaviour it exercises; any unexpected call to an unimplemented
// method panics with a nil-pointer dereference, which surfaces the mistake
// loudly rather than silently passing.

// ---------- fakeRequestRepo ----------

type fakeRequestRepo struct {
	repositories.IRequestRepo

	getRequestByID func(ctx context.Context, id int) (models.RequestDTORead, error)
}

func (f *fakeRequestRepo) GetRequestByID(ctx context.Context, id int) (models.RequestDTORead, error) {
	return f.getRequestByID(ctx, id)
}

// ---------- fakeCommentRepo ----------

type fakeCommentRepo struct {
	repositories.IRequestCommentRepo

	getLastCommentFromUser func(ctx context.Context, userID int) (models.RequestComment, error)
	addComment             func(ctx context.Context, comment models.RequestComment) (int, error)

	// addedComment captures the last comment passed to AddComment so tests can
	// assert on the values the service populated (UserID, RequestID, timestamps).
	addedComment *models.RequestComment
}

func (f *fakeCommentRepo) GetLastCommentFromUser(ctx context.Context, userID int) (models.RequestComment, error) {
	return f.getLastCommentFromUser(ctx, userID)
}

func (f *fakeCommentRepo) AddComment(ctx context.Context, comment models.RequestComment) (int, error) {
	f.addedComment = &comment
	return f.addComment(ctx, comment)
}

// ---------- fakeInvitationRepo ----------

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

// ---------- fakeDepartmentRepo ----------

type fakeDepartmentRepo struct {
	repositories.IDepartmentRepo

	getDepartmentByID func(ctx context.Context, id int) (models.DepartmentDTORead, error)
}

func (f *fakeDepartmentRepo) GetDepartmentByID(ctx context.Context, id int) (models.DepartmentDTORead, error) {
	return f.getDepartmentByID(ctx, id)
}

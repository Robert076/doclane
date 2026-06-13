package services

import (
	"context"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/Robert076/doclane/backend/models"
	"github.com/Robert076/doclane/backend/types"
	"github.com/Robert076/doclane/backend/types/errors"
)

// discardLogger returns a logger that throws away output, keeping test runs quiet.
func discardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func adminClaims() types.CallerContext {
	return types.CallerContext{UserID: 1, Role: types.RoleAdmin}
}

func memberClaims(userID int, departmentID int) types.CallerContext {
	return types.CallerContext{UserID: userID, Role: types.RoleMember, DepartmentID: &departmentID}
}

// openRequestInDepartment builds a request that is neither closed nor cancelled.
func openRequestInDepartment(departmentID, assignee int) models.RequestDTORead {
	var req models.RequestDTORead
	req.DepartmentID = departmentID
	req.Assignee = assignee
	return req
}

// ---------- validateComment ----------

func TestValidateComment_ValidCommentPasses(t *testing.T) {
	s := &RequestCommentService{logger: discardLogger()}

	err := s.validateComment(models.RequestComment{Comment: "This looks fine."})

	if err != nil {
		t.Errorf("expected no error for a valid comment, got %v", err)
	}
}

func TestValidateComment_TooShortFails(t *testing.T) {
	s := &RequestCommentService{logger: discardLogger()}

	err := s.validateComment(models.RequestComment{Comment: "hi"})

	if err == nil {
		t.Error("expected an error for a comment shorter than 3 characters")
	}
}

func TestValidateComment_WhitespaceOnlyFails(t *testing.T) {
	s := &RequestCommentService{logger: discardLogger()}

	err := s.validateComment(models.RequestComment{Comment: "     "})

	if err == nil {
		t.Error("expected an error for a whitespace-only comment")
	}
}

func TestValidateComment_TooLongFails(t *testing.T) {
	s := &RequestCommentService{logger: discardLogger()}
	long := make([]byte, 201)
	for i := range long {
		long[i] = 'a'
	}

	err := s.validateComment(models.RequestComment{Comment: string(long)})

	if err == nil {
		t.Error("expected an error for a comment longer than 200 characters")
	}
}

// ---------- checkUserIsNotSpamming ----------

func TestCheckUserIsNotSpamming_RecentCommentIsBlocked(t *testing.T) {
	commentRepo := &fakeCommentRepo{
		getLastCommentFromUser: func(ctx context.Context, userID int) (models.RequestComment, error) {
			return models.RequestComment{CreatedAt: time.Now().UTC()}, nil
		},
	}
	s := &RequestCommentService{commentRepo: commentRepo, logger: discardLogger()}

	err := s.checkUserIsNotSpamming(context.Background(), 1)

	if err == nil {
		t.Error("expected a too-many-requests error for a comment posted seconds ago")
	}
	if !errors.IsTooManyRequests(err) {
		t.Errorf("expected ErrTooManyRequests, got %T", err)
	}
}

func TestCheckUserIsNotSpamming_OldCommentIsAllowed(t *testing.T) {
	commentRepo := &fakeCommentRepo{
		getLastCommentFromUser: func(ctx context.Context, userID int) (models.RequestComment, error) {
			return models.RequestComment{CreatedAt: time.Now().UTC().Add(-time.Minute)}, nil
		},
	}
	s := &RequestCommentService{commentRepo: commentRepo, logger: discardLogger()}

	err := s.checkUserIsNotSpamming(context.Background(), 1)

	if err != nil {
		t.Errorf("expected no error when the last comment is over 30s old, got %v", err)
	}
}

func TestCheckUserIsNotSpamming_NoPreviousCommentIsAllowed(t *testing.T) {
	commentRepo := &fakeCommentRepo{
		getLastCommentFromUser: func(ctx context.Context, userID int) (models.RequestComment, error) {
			return models.RequestComment{}, errors.ErrNotFound{Msg: "no comments"}
		},
	}
	s := &RequestCommentService{commentRepo: commentRepo, logger: discardLogger()}

	err := s.checkUserIsNotSpamming(context.Background(), 1)

	if err != nil {
		t.Errorf("expected no error when the user has never commented, got %v", err)
	}
}

// ---------- checkUserIsParticipantOfRequest ----------

func TestCheckUserIsParticipant_AdminAlwaysAllowed(t *testing.T) {
	requestRepo := &fakeRequestRepo{
		getRequestByID: func(ctx context.Context, id int) (models.RequestDTORead, error) {
			return openRequestInDepartment(5, 99), nil
		},
	}
	s := &RequestCommentService{requestRepo: requestRepo, logger: discardLogger()}

	_, err := s.checkUserIsParticipantOfRequest(context.Background(), adminClaims(), 1)

	if err != nil {
		t.Errorf("expected admin to be allowed, got %v", err)
	}
}

func TestCheckUserIsParticipant_SameDepartmentAllowed(t *testing.T) {
	requestRepo := &fakeRequestRepo{
		getRequestByID: func(ctx context.Context, id int) (models.RequestDTORead, error) {
			return openRequestInDepartment(5, 99), nil
		},
	}
	s := &RequestCommentService{requestRepo: requestRepo, logger: discardLogger()}

	_, err := s.checkUserIsParticipantOfRequest(context.Background(), memberClaims(2, 5), 1)

	if err != nil {
		t.Errorf("expected a member of the same department to be allowed, got %v", err)
	}
}

func TestCheckUserIsParticipant_AssigneeAllowed(t *testing.T) {
	requestRepo := &fakeRequestRepo{
		getRequestByID: func(ctx context.Context, id int) (models.RequestDTORead, error) {
			return openRequestInDepartment(5, 2), nil
		},
	}
	s := &RequestCommentService{requestRepo: requestRepo, logger: discardLogger()}

	// caller is in a different department (7) but is the assignee (user 2)
	_, err := s.checkUserIsParticipantOfRequest(context.Background(), memberClaims(2, 7), 1)

	if err != nil {
		t.Errorf("expected the assignee to be allowed, got %v", err)
	}
}

func TestCheckUserIsParticipant_OutsiderForbidden(t *testing.T) {
	requestRepo := &fakeRequestRepo{
		getRequestByID: func(ctx context.Context, id int) (models.RequestDTORead, error) {
			return openRequestInDepartment(5, 99), nil
		},
	}
	s := &RequestCommentService{requestRepo: requestRepo, logger: discardLogger()}

	_, err := s.checkUserIsParticipantOfRequest(context.Background(), memberClaims(2, 7), 1)

	if err == nil {
		t.Error("expected a forbidden error for an unrelated member")
	}
	if !errors.IsForbidden(err) {
		t.Errorf("expected ErrForbidden, got %T", err)
	}
}

// ---------- AddComment ----------

func TestAddComment_HappyPathPopulatesFieldsAndReturnsID(t *testing.T) {
	requestRepo := &fakeRequestRepo{
		getRequestByID: func(ctx context.Context, id int) (models.RequestDTORead, error) {
			return openRequestInDepartment(5, 99), nil
		},
	}
	commentRepo := &fakeCommentRepo{
		getLastCommentFromUser: func(ctx context.Context, userID int) (models.RequestComment, error) {
			return models.RequestComment{}, errors.ErrNotFound{Msg: "none"}
		},
		addComment: func(ctx context.Context, comment models.RequestComment) (int, error) {
			return 42, nil
		},
	}
	s := &RequestCommentService{commentRepo: commentRepo, requestRepo: requestRepo, logger: discardLogger()}

	id, err := s.AddComment(context.Background(), adminClaims(), 1, models.RequestComment{Comment: "A valid comment."})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if id == nil || *id != 42 {
		t.Fatalf("expected returned id 42, got %v", id)
	}
	if commentRepo.addedComment == nil {
		t.Fatal("expected the comment to be persisted")
	}
	if commentRepo.addedComment.RequestID != 1 {
		t.Errorf("expected RequestID to be set to 1, got %d", commentRepo.addedComment.RequestID)
	}
	if commentRepo.addedComment.UserID != 1 {
		t.Errorf("expected UserID to be set from claims (1), got %d", commentRepo.addedComment.UserID)
	}
	if commentRepo.addedComment.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
}

func TestAddComment_InvalidCommentNotPersisted(t *testing.T) {
	commentRepo := &fakeCommentRepo{
		addComment: func(ctx context.Context, comment models.RequestComment) (int, error) {
			t.Fatal("AddComment should not be called for an invalid comment")
			return 0, nil
		},
	}
	s := &RequestCommentService{commentRepo: commentRepo, logger: discardLogger()}

	_, err := s.AddComment(context.Background(), adminClaims(), 1, models.RequestComment{Comment: "no"})

	if err == nil {
		t.Error("expected a validation error for too short a comment")
	}
}

func TestAddComment_ClosedRequestRejected(t *testing.T) {
	requestRepo := &fakeRequestRepo{
		getRequestByID: func(ctx context.Context, id int) (models.RequestDTORead, error) {
			req := openRequestInDepartment(5, 99)
			req.IsClosed = true
			return req, nil
		},
	}
	commentRepo := &fakeCommentRepo{
		getLastCommentFromUser: func(ctx context.Context, userID int) (models.RequestComment, error) {
			return models.RequestComment{}, errors.ErrNotFound{Msg: "none"}
		},
		addComment: func(ctx context.Context, comment models.RequestComment) (int, error) {
			t.Fatal("AddComment should not be called for a closed request")
			return 0, nil
		},
	}
	s := &RequestCommentService{commentRepo: commentRepo, requestRepo: requestRepo, logger: discardLogger()}

	_, err := s.AddComment(context.Background(), adminClaims(), 1, models.RequestComment{Comment: "A valid comment."})

	if err == nil {
		t.Error("expected an error when commenting on a closed request")
	}
}

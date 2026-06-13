//go:build integration

package repositories

import (
	"context"
	"testing"
	"time"

	"github.com/Robert076/doclane/backend/models"
	apierrors "github.com/Robert076/doclane/backend/types/errors"
)

func TestRequestCommentRepo_AddAndGetRoundTrip(t *testing.T) {
	resetDB(t)
	deptID := seedDepartment(t, "Registry")
	userID := seedUser(t, "sub-1", "ana@example.com", "Ana", "Pop", "member", deptID)
	reqID := seedRequest(t, "Birth certificate", userID, deptID)

	repo := NewRequestCommentRepo(testDB)
	now := time.Now().UTC().Truncate(time.Millisecond)

	id, err := repo.AddComment(context.Background(), models.RequestComment{
		RequestID: reqID,
		UserID:    userID,
		Comment:   "Please re-upload page 2.",
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		t.Fatalf("AddComment: %v", err)
	}
	if id == 0 {
		t.Fatal("expected a non-zero generated id")
	}

	got, err := repo.GetCommentByID(context.Background(), id)
	if err != nil {
		t.Fatalf("GetCommentByID: %v", err)
	}

	if got.ID != id {
		t.Errorf("expected id %d, got %d", id, got.ID)
	}
	if got.Comment != "Please re-upload page 2." {
		t.Errorf("comment text mismatch: got %q", got.Comment)
	}
	if got.UserFirstName != "Ana" || got.UserLastName != "Pop" {
		t.Errorf("expected author Ana Pop from the users join, got %q %q", got.UserFirstName, got.UserLastName)
	}
}

func TestRequestCommentRepo_GetCommentByID_NotFound(t *testing.T) {
	resetDB(t)

	repo := NewRequestCommentRepo(testDB)

	_, err := repo.GetCommentByID(context.Background(), 999999)

	if err == nil {
		t.Fatal("expected an error for a non-existent comment id")
	}
	if !apierrors.IsNotFound(err) {
		t.Errorf("expected ErrNotFound, got %T (%v)", err, err)
	}
}

func TestRequestCommentRepo_GetCommentsByRequestID(t *testing.T) {
	resetDB(t)
	deptID := seedDepartment(t, "Registry")
	userID := seedUser(t, "sub-1", "ana@example.com", "Ana", "Pop", "member", deptID)
	reqID := seedRequest(t, "Birth certificate", userID, deptID)
	otherReqID := seedRequest(t, "Marriage certificate", userID, deptID)

	repo := NewRequestCommentRepo(testDB)
	now := time.Now().UTC()

	for i := 0; i < 3; i++ {
		if _, err := repo.AddComment(context.Background(), models.RequestComment{
			RequestID: reqID, UserID: userID, Comment: "c", CreatedAt: now, UpdatedAt: now,
		}); err != nil {
			t.Fatalf("seeding comment: %v", err)
		}
	}
	if _, err := repo.AddComment(context.Background(), models.RequestComment{
		RequestID: otherReqID, UserID: userID, Comment: "other", CreatedAt: now, UpdatedAt: now,
	}); err != nil {
		t.Fatalf("seeding other comment: %v", err)
	}

	comments, err := repo.GetCommentsByRequestID(context.Background(), reqID)
	if err != nil {
		t.Fatalf("GetCommentsByRequestID: %v", err)
	}
	if len(comments) != 3 {
		t.Errorf("expected exactly 3 comments scoped to the request, got %d", len(comments))
	}
}

func TestRequestCommentRepo_GetLastCommentFromUser(t *testing.T) {
	resetDB(t)
	deptID := seedDepartment(t, "Registry")
	userID := seedUser(t, "sub-1", "ana@example.com", "Ana", "Pop", "member", deptID)
	reqID := seedRequest(t, "Birth certificate", userID, deptID)

	repo := NewRequestCommentRepo(testDB)
	older := time.Now().UTC().Add(-time.Hour)
	newer := time.Now().UTC()

	if _, err := repo.AddComment(context.Background(), models.RequestComment{
		RequestID: reqID, UserID: userID, Comment: "older", CreatedAt: older, UpdatedAt: older,
	}); err != nil {
		t.Fatalf("seeding older comment: %v", err)
	}
	if _, err := repo.AddComment(context.Background(), models.RequestComment{
		RequestID: reqID, UserID: userID, Comment: "newer", CreatedAt: newer, UpdatedAt: newer,
	}); err != nil {
		t.Fatalf("seeding newer comment: %v", err)
	}

	last, err := repo.GetLastCommentFromUser(context.Background(), userID)
	if err != nil {
		t.Fatalf("GetLastCommentFromUser: %v", err)
	}
	if last.Comment != "newer" {
		t.Errorf("expected the most recent comment to be returned, got %q", last.Comment)
	}
}

//go:build integration

package repositories

import (
	"context"
	"database/sql"
	"testing"

	"github.com/Robert076/doclane/backend/models"
)

// buildRequest assembles a minimal valid Request for insertion.
func buildRequest(title string, assignee, departmentID int) models.Request {
	var req models.Request
	req.Title = title
	req.Assignee = assignee
	req.DepartmentID = departmentID
	return req
}

func TestRequestRepo_AddRequestWithTx_CommitCreatesRequestAndDocs(t *testing.T) {
	resetDB(t)
	deptID := seedDepartment(t, "Registry")
	userID := seedUser(t, "sub-1", "ana@example.com", "Ana", "Pop", "member", deptID)

	requestRepo := NewRequestRepo(testDB)
	expectedDocRepo := NewExpectedDocRepo(testDB)
	mgr := NewTxManager(testDB)

	var requestID int
	err := mgr.WithTx(context.Background(), func(tx *sql.Tx) error {
		id, err := requestRepo.AddRequestWithTx(context.Background(), buildRequest("ID card", userID, deptID), tx)
		if err != nil {
			return err
		}
		requestID = id

		for _, title := range []string{"Old ID scan", "Proof of address"} {
			ed := models.ExpectedDocument{RequestID: id, Title: title, Status: "pending"}
			if err := expectedDocRepo.AddExpectedDocumentToRequestWithTx(context.Background(), tx, ed); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected the transaction to commit, got %v", err)
	}

	// The request is readable...
	got, err := requestRepo.GetRequestByID(context.Background(), requestID)
	if err != nil {
		t.Fatalf("GetRequestByID: %v", err)
	}
	if got.Title != "ID card" {
		t.Errorf("expected title 'ID card', got %q", got.Title)
	}

	// ...and both expected documents were persisted atomically with it.
	docs, err := expectedDocRepo.GetExpectedDocumentsByRequest(context.Background(), requestID)
	if err != nil {
		t.Fatalf("GetExpectedDocumentsByRequest: %v", err)
	}
	if len(docs) != 2 {
		t.Errorf("expected 2 expected documents, got %d", len(docs))
	}
}

func TestRequestRepo_AddRequestWithTx_RollbackLeavesNoOrphans(t *testing.T) {
	resetDB(t)
	deptID := seedDepartment(t, "Registry")
	userID := seedUser(t, "sub-1", "ana@example.com", "Ana", "Pop", "member", deptID)

	requestRepo := NewRequestRepo(testDB)
	expectedDocRepo := NewExpectedDocRepo(testDB)
	mgr := NewTxManager(testDB)

	err := mgr.WithTx(context.Background(), func(tx *sql.Tx) error {
		id, err := requestRepo.AddRequestWithTx(context.Background(), buildRequest("ID card", userID, deptID), tx)
		if err != nil {
			return err
		}
		// A document with an invalid status violates the CHECK constraint, which
		// must roll back the request inserted moments earlier in the same tx.
		ed := models.ExpectedDocument{RequestID: id, Title: "bad", Status: "not-a-valid-status"}
		return expectedDocRepo.AddExpectedDocumentToRequestWithTx(context.Background(), tx, ed)
	})
	if err == nil {
		t.Fatal("expected the invalid status to fail the transaction")
	}

	var requestCount int
	if err := testDB.QueryRowContext(context.Background(),
		`SELECT COUNT(*) FROM document_requests`).Scan(&requestCount); err != nil {
		t.Fatalf("counting requests: %v", err)
	}
	if requestCount != 0 {
		t.Errorf("expected no request rows after rollback, got %d (orphaned request)", requestCount)
	}
}

func TestRequestRepo_ClaimAndUnclaim(t *testing.T) {
	resetDB(t)
	deptID := seedDepartment(t, "Registry")
	userID := seedUser(t, "sub-1", "ana@example.com", "Ana", "Pop", "member", deptID)
	reqID := seedRequest(t, "ID card", userID, deptID)

	repo := NewRequestRepo(testDB)

	if err := repo.ClaimRequest(context.Background(), reqID, userID); err != nil {
		t.Fatalf("ClaimRequest: %v", err)
	}
	claimed, err := repo.GetRequestByID(context.Background(), reqID)
	if err != nil {
		t.Fatalf("GetRequestByID after claim: %v", err)
	}
	if claimed.ClaimedBy == nil || *claimed.ClaimedBy != userID {
		t.Errorf("expected claimed_by to be %d, got %v", userID, claimed.ClaimedBy)
	}

	if err := repo.UnclaimRequest(context.Background(), reqID); err != nil {
		t.Fatalf("UnclaimRequest: %v", err)
	}
	unclaimed, err := repo.GetRequestByID(context.Background(), reqID)
	if err != nil {
		t.Fatalf("GetRequestByID after unclaim: %v", err)
	}
	if unclaimed.ClaimedBy != nil {
		t.Errorf("expected claimed_by to be NULL after unclaim, got %v", *unclaimed.ClaimedBy)
	}
}

func TestRequestRepo_CloseAndReopen(t *testing.T) {
	resetDB(t)
	deptID := seedDepartment(t, "Registry")
	userID := seedUser(t, "sub-1", "ana@example.com", "Ana", "Pop", "member", deptID)
	reqID := seedRequest(t, "ID card", userID, deptID)

	repo := NewRequestRepo(testDB)

	if err := repo.CloseRequest(context.Background(), reqID); err != nil {
		t.Fatalf("CloseRequest: %v", err)
	}
	closed, err := repo.GetRequestByID(context.Background(), reqID)
	if err != nil {
		t.Fatalf("GetRequestByID after close: %v", err)
	}
	if !closed.IsClosed {
		t.Error("expected request to be closed")
	}

	if err := repo.ReopenRequest(context.Background(), reqID); err != nil {
		t.Fatalf("ReopenRequest: %v", err)
	}
	reopened, err := repo.GetRequestByID(context.Background(), reqID)
	if err != nil {
		t.Fatalf("GetRequestByID after reopen: %v", err)
	}
	if reopened.IsClosed {
		t.Error("expected request to be reopened (is_closed = false)")
	}
}

func TestRequestRepo_Cancel(t *testing.T) {
	resetDB(t)
	deptID := seedDepartment(t, "Registry")
	userID := seedUser(t, "sub-1", "ana@example.com", "Ana", "Pop", "member", deptID)
	reqID := seedRequest(t, "ID card", userID, deptID)

	repo := NewRequestRepo(testDB)

	if err := repo.CancelRequest(context.Background(), reqID); err != nil {
		t.Fatalf("CancelRequest: %v", err)
	}
	cancelled, err := repo.GetRequestByID(context.Background(), reqID)
	if err != nil {
		t.Fatalf("GetRequestByID after cancel: %v", err)
	}
	if !cancelled.IsCancelled {
		t.Error("expected request to be cancelled")
	}
}

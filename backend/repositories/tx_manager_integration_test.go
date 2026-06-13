//go:build integration

package repositories

import (
	"context"
	"database/sql"
	"errors"
	"testing"
)

func countComments(t *testing.T, requestID int) int {
	t.Helper()
	var n int
	err := testDB.QueryRowContext(context.Background(),
		`SELECT COUNT(*) FROM document_comments WHERE request_id = $1`, requestID,
	).Scan(&n)
	if err != nil {
		t.Fatalf("countComments: %v", err)
	}
	return n
}

func TestTxManager_CommitPersistsAllWrites(t *testing.T) {
	resetDB(t)
	deptID := seedDepartment(t, "Registry")
	userID := seedUser(t, "sub-1", "user@example.com", "Ana", "Pop", "member", deptID)
	reqID := seedRequest(t, "Birth certificate", userID, deptID)

	mgr := NewTxManager(testDB)

	err := mgr.WithTx(context.Background(), func(tx *sql.Tx) error {
		for i := 0; i < 2; i++ {
			if _, err := tx.ExecContext(context.Background(),
				`INSERT INTO document_comments (request_id, user_id, comment) VALUES ($1, $2, $3)`,
				reqID, userID, "a committed comment",
			); err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		t.Fatalf("expected commit to succeed, got %v", err)
	}
	if got := countComments(t, reqID); got != 2 {
		t.Errorf("expected 2 persisted comments after commit, got %d", got)
	}
}

func TestTxManager_RollsBackAllWritesOnError(t *testing.T) {
	resetDB(t)
	deptID := seedDepartment(t, "Registry")
	userID := seedUser(t, "sub-1", "user@example.com", "Ana", "Pop", "member", deptID)
	reqID := seedRequest(t, "Birth certificate", userID, deptID)

	mgr := NewTxManager(testDB)
	sentinel := errors.New("boom halfway through")

	err := mgr.WithTx(context.Background(), func(tx *sql.Tx) error {
		if _, err := tx.ExecContext(context.Background(),
			`INSERT INTO document_comments (request_id, user_id, comment) VALUES ($1, $2, $3)`,
			reqID, userID, "this should be rolled back",
		); err != nil {
			return err
		}
		return sentinel
	})

	if !errors.Is(err, sentinel) {
		t.Fatalf("expected the sentinel error to propagate, got %v", err)
	}
	if got := countComments(t, reqID); got != 0 {
		t.Errorf("expected 0 comments after rollback, got %d", got)
	}
}

func TestTxManager_RollsBackOnConstraintViolation(t *testing.T) {
	resetDB(t)
	deptID := seedDepartment(t, "Registry")
	userID := seedUser(t, "sub-1", "user@example.com", "Ana", "Pop", "member", deptID)
	reqID := seedRequest(t, "Birth certificate", userID, deptID)

	mgr := NewTxManager(testDB)

	err := mgr.WithTx(context.Background(), func(tx *sql.Tx) error {
		if _, err := tx.ExecContext(context.Background(),
			`INSERT INTO document_comments (request_id, user_id, comment) VALUES ($1, $2, $3)`,
			reqID, userID, "valid",
		); err != nil {
			return err
		}
		_, err := tx.ExecContext(context.Background(),
			`INSERT INTO document_comments (request_id, user_id, comment) VALUES ($1, $2, $3)`,
			reqID, 999999, "invalid user",
		)
		return err
	})

	if err == nil {
		t.Fatal("expected a constraint-violation error")
	}
	if got := countComments(t, reqID); got != 0 {
		t.Errorf("expected 0 comments after a failed transaction, got %d", got)
	}
}

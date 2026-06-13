//go:build integration

package repositories

import (
	"context"
	"testing"
	"time"
)

func seedDepartment(t *testing.T, name string) int {
	t.Helper()
	var id int
	err := testDB.QueryRowContext(context.Background(),
		`INSERT INTO departments (name) VALUES ($1) RETURNING id`,
		name,
	).Scan(&id)
	if err != nil {
		t.Fatalf("seedDepartment(%q): %v", name, err)
	}
	return id
}

func seedUser(t *testing.T, cognitoSub, email, firstName, lastName, role string, departmentID int) int {
	t.Helper()
	var deptArg any
	if departmentID != 0 {
		deptArg = departmentID
	}
	var id int
	err := testDB.QueryRowContext(context.Background(),
		`INSERT INTO users (cognito_sub, email, first_name, last_name, role, department_id)
		 VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		cognitoSub, email, firstName, lastName, role, deptArg,
	).Scan(&id)
	if err != nil {
		t.Fatalf("seedUser(%q): %v", email, err)
	}
	return id
}

func seedRequest(t *testing.T, title string, assignee, departmentID int) int {
	t.Helper()
	var id int
	err := testDB.QueryRowContext(context.Background(),
		`INSERT INTO document_requests (title, assignee, department_id)
		 VALUES ($1, $2, $3) RETURNING id`,
		title, assignee, departmentID,
	).Scan(&id)
	if err != nil {
		t.Fatalf("seedRequest(%q): %v", title, err)
	}
	return id
}

func seedInvitationCode(t *testing.T, code string, createdBy, departmentID int, expiresAt *time.Time) int {
	t.Helper()
	var id int
	err := testDB.QueryRowContext(context.Background(),
		`INSERT INTO invitation_codes (code, created_by, department_id, expires_at)
		 VALUES ($1, $2, $3, $4) RETURNING id`,
		code, createdBy, departmentID, expiresAt,
	).Scan(&id)
	if err != nil {
		t.Fatalf("seedInvitationCode(%q): %v", code, err)
	}
	return id
}

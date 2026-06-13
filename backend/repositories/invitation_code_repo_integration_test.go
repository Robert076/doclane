//go:build integration

package repositories

import (
	"context"
	"testing"
	"time"

	apierrors "github.com/Robert076/doclane/backend/types/errors"
)

func TestInvitationCodeRepo_CreateAndGetByCode(t *testing.T) {
	resetDB(t)
	deptID := seedDepartment(t, "Registry")
	adminID := seedUser(t, "sub-admin", "admin@example.com", "Admin", "User", "admin", deptID)

	repo := NewInvitationCodeRepo(testDB)

	if err := repo.CreateInvitationCode(context.Background(), deptID, "ABCD-EF01-2345", adminID, nil); err != nil {
		t.Fatalf("CreateInvitationCode: %v", err)
	}

	got, err := repo.GetInvitationCodeByCode(context.Background(), "ABCD-EF01-2345")
	if err != nil {
		t.Fatalf("GetInvitationCodeByCode: %v", err)
	}
	if got.Code != "ABCD-EF01-2345" {
		t.Errorf("expected the created code back, got %q", got.Code)
	}
	if got.DepartmentName != "Registry" {
		t.Errorf("expected department name 'Registry', got %q", got.DepartmentName)
	}
}

func TestInvitationCodeRepo_GetByCode_NotFound(t *testing.T) {
	resetDB(t)

	repo := NewInvitationCodeRepo(testDB)

	_, err := repo.GetInvitationCodeByCode(context.Background(), "NOPE-NOPE-NOPE")

	if !apierrors.IsNotFound(err) {
		t.Errorf("expected ErrNotFound for an unknown code, got %T (%v)", err, err)
	}
}

func TestInvitationCodeRepo_GetByCreator_ExcludesUsedCodes(t *testing.T) {
	resetDB(t)
	deptID := seedDepartment(t, "Registry")
	adminID := seedUser(t, "sub-admin", "admin@example.com", "Admin", "User", "admin", deptID)

	repo := NewInvitationCodeRepo(testDB)

	if err := repo.CreateInvitationCode(context.Background(), deptID, "AAAA-AAAA-AAAA", adminID, nil); err != nil {
		t.Fatalf("create unused: %v", err)
	}
	usedID := seedInvitationCode(t, "BBBB-BBBB-BBBB", adminID, deptID, nil)
	if err := repo.InvalidateCode(context.Background(), usedID, adminID); err != nil {
		t.Fatalf("InvalidateCode: %v", err)
	}

	codes, err := repo.GetInvitationCodesByCreator(context.Background(), adminID)
	if err != nil {
		t.Fatalf("GetInvitationCodesByCreator: %v", err)
	}
	if len(codes) != 1 {
		t.Fatalf("expected only the 1 unused code, got %d", len(codes))
	}
	if codes[0].Code != "AAAA-AAAA-AAAA" {
		t.Errorf("expected the unused code, got %q", codes[0].Code)
	}
}

func TestInvitationCodeRepo_InvalidateCode_IsIdempotentlyGuarded(t *testing.T) {
	resetDB(t)
	deptID := seedDepartment(t, "Registry")
	adminID := seedUser(t, "sub-admin", "admin@example.com", "Admin", "User", "admin", deptID)
	codeID := seedInvitationCode(t, "CCCC-CCCC-CCCC", adminID, deptID, nil)

	repo := NewInvitationCodeRepo(testDB)

	if err := repo.InvalidateCode(context.Background(), codeID, adminID); err != nil {
		t.Fatalf("first InvalidateCode: %v", err)
	}
	if err := repo.InvalidateCode(context.Background(), codeID, adminID); err == nil {
		t.Error("expected an error when invalidating an already-used code")
	}
}

func TestInvitationCodeRepo_DeleteCode(t *testing.T) {
	resetDB(t)
	deptID := seedDepartment(t, "Registry")
	adminID := seedUser(t, "sub-admin", "admin@example.com", "Admin", "User", "admin", deptID)
	codeID := seedInvitationCode(t, "DDDD-DDDD-DDDD", adminID, deptID, nil)

	repo := NewInvitationCodeRepo(testDB)

	if err := repo.DeleteCode(context.Background(), codeID); err != nil {
		t.Fatalf("DeleteCode: %v", err)
	}

	if err := repo.DeleteCode(context.Background(), codeID); err == nil {
		t.Error("expected an error when deleting a non-existent code")
	}
}

func TestInvitationCodeRepo_GetByID_PersistsExpiry(t *testing.T) {
	resetDB(t)
	deptID := seedDepartment(t, "Registry")
	adminID := seedUser(t, "sub-admin", "admin@example.com", "Admin", "User", "admin", deptID)

	repo := NewInvitationCodeRepo(testDB)
	expiry := time.Now().Add(72 * time.Hour).UTC().Truncate(time.Second)

	if err := repo.CreateInvitationCode(context.Background(), deptID, "EEEE-EEEE-EEEE", adminID, &expiry); err != nil {
		t.Fatalf("CreateInvitationCode: %v", err)
	}

	codes, err := repo.GetInvitationCodesByCreator(context.Background(), adminID)
	if err != nil {
		t.Fatalf("GetInvitationCodesByCreator: %v", err)
	}
	if len(codes) != 1 {
		t.Fatalf("expected 1 code, got %d", len(codes))
	}
	if codes[0].ExpiresAt == nil {
		t.Fatal("expected the expiry to be persisted")
	}
	if !codes[0].ExpiresAt.UTC().Truncate(time.Second).Equal(expiry) {
		t.Errorf("expected expiry %v, got %v", expiry, codes[0].ExpiresAt.UTC())
	}
}

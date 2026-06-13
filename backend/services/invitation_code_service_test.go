package services

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/Robert076/doclane/backend/models"
	"github.com/Robert076/doclane/backend/types/errors"
)

// existingDepartment is a department lookup that always succeeds.
func departmentExists(ctx context.Context, id int) (models.DepartmentDTORead, error) {
	return models.DepartmentDTORead{}, nil
}

// ---------- generateInvitationCode ----------

func TestGenerateInvitationCode_FormatIsCorrect(t *testing.T) {
	code, err := generateInvitationCode()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	matched := regexp.MustCompile(`^[0-9A-F]{4}-[0-9A-F]{4}-[0-9A-F]{4}$`).MatchString(code)
	if !matched {
		t.Errorf("expected code in XXXX-XXXX-XXXX uppercase hex format, got %q", code)
	}
}

func TestGenerateInvitationCode_ProducesDistinctValues(t *testing.T) {
	first, _ := generateInvitationCode()
	second, _ := generateInvitationCode()

	if first == second {
		t.Errorf("expected two generated codes to differ, both were %q", first)
	}
}

// ---------- CreateInvitationCode ----------

func TestCreateInvitationCode_NonAdminForbidden(t *testing.T) {
	s := &InvitationCodeService{logger: discardLogger()}

	_, err := s.CreateInvitationCode(context.Background(), memberClaims(2, 5), 5, 7)

	if !errors.IsForbidden(err) {
		t.Errorf("expected ErrForbidden for a non-admin, got %T", err)
	}
}

func TestCreateInvitationCode_UnknownDepartmentFails(t *testing.T) {
	deptRepo := &fakeDepartmentRepo{
		getDepartmentByID: func(ctx context.Context, id int) (models.DepartmentDTORead, error) {
			return models.DepartmentDTORead{}, errors.ErrNotFound{Msg: "nope"}
		},
	}
	s := &InvitationCodeService{departmentRepo: deptRepo, logger: discardLogger()}

	_, err := s.CreateInvitationCode(context.Background(), adminClaims(), 99, 7)

	if !errors.IsNotFound(err) {
		t.Errorf("expected ErrNotFound for an unknown department, got %T", err)
	}
}

func TestCreateInvitationCode_TooManyActiveCodesFails(t *testing.T) {
	deptRepo := &fakeDepartmentRepo{getDepartmentByID: departmentExists}
	invRepo := &fakeInvitationRepo{
		getByCreator: func(ctx context.Context, createdBy int) ([]models.InvitationCode, error) {
			return []models.InvitationCode{{}, {}, {}}, nil // already at the limit of 3
		},
	}
	s := &InvitationCodeService{invitationRepo: invRepo, departmentRepo: deptRepo, logger: discardLogger()}

	_, err := s.CreateInvitationCode(context.Background(), adminClaims(), 5, 7)

	if !errors.IsBadRequest(err) {
		t.Errorf("expected ErrBadRequest when at the 3-code limit, got %T", err)
	}
}

func TestCreateInvitationCode_HappyPathPersistsCode(t *testing.T) {
	deptRepo := &fakeDepartmentRepo{getDepartmentByID: departmentExists}
	invRepo := &fakeInvitationRepo{
		getByCreator: func(ctx context.Context, createdBy int) ([]models.InvitationCode, error) {
			return []models.InvitationCode{}, nil
		},
		createInvitationCode: func(ctx context.Context, departmentID int, code string, createdBy int, expiresAt *time.Time) error {
			return nil
		},
	}
	s := &InvitationCodeService{invitationRepo: invRepo, departmentRepo: deptRepo, logger: discardLogger()}

	code, err := s.CreateInvitationCode(context.Background(), adminClaims(), 5, 7)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if code == "" {
		t.Error("expected a non-empty code to be returned")
	}
	if invRepo.createdCode != code {
		t.Errorf("expected the persisted code (%q) to match the returned code (%q)", invRepo.createdCode, code)
	}
	if invRepo.createdExpiresAt == nil {
		t.Error("expected an expiry to be set when expiresInDays > 0")
	}
}

func TestCreateInvitationCode_NoExpiryWhenDaysZero(t *testing.T) {
	deptRepo := &fakeDepartmentRepo{getDepartmentByID: departmentExists}
	invRepo := &fakeInvitationRepo{
		getByCreator: func(ctx context.Context, createdBy int) ([]models.InvitationCode, error) {
			return []models.InvitationCode{}, nil
		},
		createInvitationCode: func(ctx context.Context, departmentID int, code string, createdBy int, expiresAt *time.Time) error {
			return nil
		},
	}
	s := &InvitationCodeService{invitationRepo: invRepo, departmentRepo: deptRepo, logger: discardLogger()}

	if _, err := s.CreateInvitationCode(context.Background(), adminClaims(), 5, 0); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if invRepo.createdExpiresAt != nil {
		t.Errorf("expected no expiry when expiresInDays is 0, got %v", invRepo.createdExpiresAt)
	}
}

// ---------- GetInvitationCodeInfo ----------

func TestGetInvitationCodeInfo_UnknownCodeIsNotFound(t *testing.T) {
	invRepo := &fakeInvitationRepo{
		getByCode: func(ctx context.Context, code string) (models.InvitationCodeReadDTO, error) {
			return models.InvitationCodeReadDTO{}, errors.ErrNotFound{Msg: "nope"}
		},
	}
	s := &InvitationCodeService{invitationRepo: invRepo, logger: discardLogger()}

	_, err := s.GetInvitationCodeInfo(context.Background(), "ABCD-EF01-2345")

	if !errors.IsNotFound(err) {
		t.Errorf("expected ErrNotFound for an unknown code, got %T", err)
	}
}

func TestGetInvitationCodeInfo_UsedCodeIsRejected(t *testing.T) {
	used := time.Now().Add(-time.Hour)
	invRepo := &fakeInvitationRepo{
		getByCode: func(ctx context.Context, code string) (models.InvitationCodeReadDTO, error) {
			var dto models.InvitationCodeReadDTO
			dto.UsedAt = &used
			return dto, nil
		},
	}
	s := &InvitationCodeService{invitationRepo: invRepo, logger: discardLogger()}

	_, err := s.GetInvitationCodeInfo(context.Background(), "ABCD-EF01-2345")

	if !errors.IsBadRequest(err) {
		t.Errorf("expected ErrBadRequest for an already-used code, got %T", err)
	}
}

func TestGetInvitationCodeInfo_ExpiredCodeIsRejected(t *testing.T) {
	expired := time.Now().Add(-time.Hour)
	invRepo := &fakeInvitationRepo{
		getByCode: func(ctx context.Context, code string) (models.InvitationCodeReadDTO, error) {
			var dto models.InvitationCodeReadDTO
			dto.ExpiresAt = &expired
			return dto, nil
		},
	}
	s := &InvitationCodeService{invitationRepo: invRepo, logger: discardLogger()}

	_, err := s.GetInvitationCodeInfo(context.Background(), "ABCD-EF01-2345")

	if !errors.IsBadRequest(err) {
		t.Errorf("expected ErrBadRequest for an expired code, got %T", err)
	}
}

func TestGetInvitationCodeInfo_ValidCodeReturned(t *testing.T) {
	future := time.Now().Add(48 * time.Hour)
	invRepo := &fakeInvitationRepo{
		getByCode: func(ctx context.Context, code string) (models.InvitationCodeReadDTO, error) {
			var dto models.InvitationCodeReadDTO
			dto.Code = code
			dto.ExpiresAt = &future
			return dto, nil
		},
	}
	s := &InvitationCodeService{invitationRepo: invRepo, logger: discardLogger()}

	got, err := s.GetInvitationCodeInfo(context.Background(), "ABCD-EF01-2345")

	if err != nil {
		t.Fatalf("expected no error for a valid code, got %v", err)
	}
	if got == nil || got.Code != "ABCD-EF01-2345" {
		t.Errorf("expected the valid code DTO to be returned, got %v", got)
	}
}

// ---------- DeleteInvitationCode ----------

func TestDeleteInvitationCode_NonAdminForbidden(t *testing.T) {
	s := &InvitationCodeService{logger: discardLogger()}

	err := s.DeleteInvitationCode(context.Background(), memberClaims(2, 5), 10)

	if !errors.IsForbidden(err) {
		t.Errorf("expected ErrForbidden for a non-admin, got %T", err)
	}
}

func TestDeleteInvitationCode_OwnedCodeIsDeleted(t *testing.T) {
	invRepo := &fakeInvitationRepo{
		getByID: func(ctx context.Context, id int) (models.InvitationCode, error) {
			return models.InvitationCode{ID: id, CreatedBy: 1}, nil
		},
		deleteCode: func(ctx context.Context, id int) error { return nil },
	}
	s := &InvitationCodeService{invitationRepo: invRepo, logger: discardLogger()}

	err := s.DeleteInvitationCode(context.Background(), adminClaims(), 10)

	if err != nil {
		t.Fatalf("expected no error deleting an owned code, got %v", err)
	}
	if len(invRepo.deletedIDs) != 1 || invRepo.deletedIDs[0] != 10 {
		t.Errorf("expected code 10 to be deleted, got %v", invRepo.deletedIDs)
	}
}

func TestDeleteInvitationCode_OtherUsersCodeForbidden(t *testing.T) {
	invRepo := &fakeInvitationRepo{
		getByID: func(ctx context.Context, id int) (models.InvitationCode, error) {
			return models.InvitationCode{ID: id, CreatedBy: 999}, nil // owned by someone else
		},
		deleteCode: func(ctx context.Context, id int) error {
			t.Fatal("DeleteCode should not be called for a code owned by another user")
			return nil
		},
	}
	s := &InvitationCodeService{invitationRepo: invRepo, logger: discardLogger()}

	err := s.DeleteInvitationCode(context.Background(), adminClaims(), 10)

	if !errors.IsForbidden(err) {
		t.Errorf("expected ErrForbidden when deleting another user's code, got %T", err)
	}
}

// ---------- deleteExpiredCodes ----------

func TestDeleteExpiredCodes_RemovesExpiredKeepsValid(t *testing.T) {
	expired := time.Now().Add(-time.Hour)
	future := time.Now().Add(time.Hour)

	invRepo := &fakeInvitationRepo{
		deleteCode: func(ctx context.Context, id int) error { return nil },
	}
	s := &InvitationCodeService{invitationRepo: invRepo, logger: discardLogger()}

	codes := []models.InvitationCode{
		{ID: 1, ExpiresAt: &expired}, // should be deleted
		{ID: 2, ExpiresAt: &future},  // should remain
		{ID: 3, ExpiresAt: nil},      // never expires, should remain
	}

	valid, err := s.deleteExpiredCodes(context.Background(), codes)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(valid) != 2 {
		t.Fatalf("expected 2 valid codes to remain, got %d", len(valid))
	}
	if len(invRepo.deletedIDs) != 1 || invRepo.deletedIDs[0] != 1 {
		t.Errorf("expected only the expired code (id 1) to be deleted, got %v", invRepo.deletedIDs)
	}
}

func TestDeleteExpiredCodes_NonAdminCannotReachThisButValidCodesPassThrough(t *testing.T) {
	// All codes valid: nothing should be deleted and all should pass through.
	future := time.Now().Add(time.Hour)
	invRepo := &fakeInvitationRepo{
		deleteCode: func(ctx context.Context, id int) error {
			t.Fatal("no code should be deleted when all are valid")
			return nil
		},
	}
	s := &InvitationCodeService{invitationRepo: invRepo, logger: discardLogger()}

	codes := []models.InvitationCode{{ID: 1, ExpiresAt: &future}, {ID: 2, ExpiresAt: &future}}

	valid, err := s.deleteExpiredCodes(context.Background(), codes)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(valid) != 2 {
		t.Errorf("expected all 2 codes to pass through, got %d", len(valid))
	}
}

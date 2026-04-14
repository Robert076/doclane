package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Robert076/doclane/backend/models"
	"github.com/Robert076/doclane/backend/types/errors"
)

type InvitationCodeRepo struct {
	db *sql.DB
}

func NewInvitationCodeRepo(db *sql.DB) *InvitationCodeRepo {
	return &InvitationCodeRepo{db: db}
}

func (r *InvitationCodeRepo) GetInvitationCodeByCode(ctx context.Context, code string) (models.InvitationCodeReadDTO, error) {
	query := `
		SELECT ic.id, ic.department_id, ic.code, ic.created_by, ic.used_at, ic.expires_at, ic.created_at, d.name
		FROM invitation_codes ic
		JOIN departments d ON d.id = ic.department_id
		WHERE ic.code = $1
	`
	var dto models.InvitationCodeReadDTO
	err := r.db.QueryRowContext(ctx, query, code).Scan(
		&dto.ID,
		&dto.DepartmentID,
		&dto.Code,
		&dto.CreatedBy,
		&dto.UsedAt,
		&dto.ExpiresAt,
		&dto.CreatedAt,
		&dto.DepartmentName,
	)
	if err == sql.ErrNoRows {
		return models.InvitationCodeReadDTO{}, errors.ErrNotFound{Msg: "Invitation code not found."}
	}
	return dto, err
}

func (r *InvitationCodeRepo) GetInvitationCodesByCreator(ctx context.Context, createdBy int) ([]models.InvitationCode, error) {
	query := `
		SELECT id, department_id, code, created_by, used_at, expires_at, created_at
		FROM invitation_codes
		WHERE created_by = $1 AND used_at IS NULL
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, createdBy)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var codes []models.InvitationCode
	for rows.Next() {
		var code models.InvitationCode
		if err := rows.Scan(
			&code.ID,
			&code.DepartmentID,
			&code.Code,
			&code.CreatedBy,
			&code.UsedAt,
			&code.ExpiresAt,
			&code.CreatedAt,
		); err != nil {
			return nil, err
		}
		codes = append(codes, code)
	}
	return codes, rows.Err()
}

func (r *InvitationCodeRepo) GetInvitationCodeByID(ctx context.Context, id int) (models.InvitationCode, error) {
	query := `
		SELECT id, department_id, code, created_by, used_at, expires_at, created_at
		FROM invitation_codes
		WHERE id = $1
	`
	var invCode models.InvitationCode
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&invCode.ID,
		&invCode.DepartmentID,
		&invCode.Code,
		&invCode.CreatedBy,
		&invCode.UsedAt,
		&invCode.ExpiresAt,
		&invCode.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return models.InvitationCode{}, errors.ErrNotFound{Msg: "Invitation code not found."}
	}
	return invCode, err
}

func (r *InvitationCodeRepo) GetInvitationCodesByDepartment(ctx context.Context, departmentID int) ([]models.InvitationCode, error) {
	query := `
		SELECT id, code, created_by, department_id, used_at, expires_at, created_at
		FROM invitation_codes
		WHERE department_id = $1 AND used_at IS NULL
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, departmentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	codes := []models.InvitationCode{}
	for rows.Next() {
		var code models.InvitationCode
		if err := rows.Scan(
			&code.ID,
			&code.Code,
			&code.CreatedBy,
			&code.DepartmentID,
			&code.UsedAt,
			&code.ExpiresAt,
			&code.CreatedAt,
		); err != nil {
			return nil, err
		}
		codes = append(codes, code)
	}
	return codes, rows.Err()
}

func (r *InvitationCodeRepo) CreateInvitationCode(ctx context.Context, departmentID int, code string, createdBy int, expiresAt *time.Time) error {
	query := `
		INSERT INTO invitation_codes (department_id, code, created_by, expires_at)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.ExecContext(ctx, query, departmentID, code, createdBy, expiresAt)
	return err
}

func (r *InvitationCodeRepo) InvalidateCode(ctx context.Context, id int) error {
	query := `UPDATE invitation_codes SET used_at = NOW() WHERE id = $1 AND used_at IS NULL`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("invitation code not found or already used")
	}
	return nil
}

func (r *InvitationCodeRepo) DeleteCode(ctx context.Context, id int) error {
	query := `DELETE FROM invitation_codes WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("invitation code not found")
	}
	return nil
}

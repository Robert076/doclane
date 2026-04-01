package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Robert076/doclane/backend/models"
)

type InvitationCodeRepo struct {
	db *sql.DB
}

func NewInvitationCodeRepo(db *sql.DB) *InvitationCodeRepo {
	return &InvitationCodeRepo{db: db}
}

func (r *InvitationCodeRepo) GetInvitationCodeByCode(ctx context.Context, code string) (models.InvitationCode, error) {
	query := `
		SELECT id, code, created_by, used_at, expires_at, created_at
		FROM invitation_codes
		WHERE code = $1
	`
	var invCode models.InvitationCode
	err := r.db.QueryRowContext(ctx, query, code).Scan(
		&invCode.ID,
		&invCode.Code,
		&invCode.CreatedBy,
		&invCode.UsedAt,
		&invCode.ExpiresAt,
		&invCode.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return models.InvitationCode{}, fmt.Errorf("invitation code not found")
	}
	return invCode, err
}

func (r *InvitationCodeRepo) GetInvitationCodesByCreator(ctx context.Context, createdBy int) ([]models.InvitationCode, error) {
	query := `
		SELECT id, code, created_by, used_at, expires_at, created_at
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
		SELECT id, code, created_by, used_at, expires_at, created_at
		FROM invitation_codes
		WHERE id = $1
	`
	var invCode models.InvitationCode
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&invCode.ID,
		&invCode.Code,
		&invCode.CreatedBy,
		&invCode.UsedAt,
		&invCode.ExpiresAt,
		&invCode.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return models.InvitationCode{}, fmt.Errorf("invitation code not found")
	}
	return invCode, err
}

func (r *InvitationCodeRepo) CreateInvitationCode(ctx context.Context, code string, createdBy int, expiresAt *time.Time) error {
	query := `
		INSERT INTO invitation_codes (code, created_by, expires_at)
		VALUES ($1, $2, $3)
	`
	_, err := r.db.ExecContext(ctx, query, code, createdBy, expiresAt)
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

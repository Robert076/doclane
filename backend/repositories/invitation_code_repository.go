package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Robert076/doclane/backend/models"
)

type InvitationCodeRepository struct {
	db *sql.DB
}

func NewInvitationCodeRepository(db *sql.DB) *InvitationCodeRepository {
	return &InvitationCodeRepository{db: db}
}

func (r *InvitationCodeRepository) GetInvitationCodeByCode(
	ctx context.Context,
	code string,
) (models.InvitationCode, error) {
	query := `
        SELECT id, code, professional_id, used_at, expires_at, created_at
        FROM invitation_codes
        WHERE code = $1
    `

	var invCode models.InvitationCode
	err := r.db.QueryRowContext(ctx, query, code).Scan(
		&invCode.ID,
		&invCode.Code,
		&invCode.ProfessionalID,
		&invCode.UsedAt,
		&invCode.ExpiresAt,
		&invCode.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return models.InvitationCode{}, fmt.Errorf("invitation code not found")
	}
	if err != nil {
		return models.InvitationCode{}, err
	}

	return invCode, nil
}

func (r *InvitationCodeRepository) GetInvitationCodesByProfessional(
	ctx context.Context,
	professionalID int,
) ([]models.InvitationCode, error) {
	query := `
        SELECT id, code, professional_id, used_at, expires_at, created_at
        FROM invitation_codes
        WHERE professional_id = $1 AND used_at IS NULL
        ORDER BY created_at DESC
    `

	rows, err := r.db.QueryContext(ctx, query, professionalID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var codes []models.InvitationCode
	for rows.Next() {
		var code models.InvitationCode
		err := rows.Scan(
			&code.ID,
			&code.Code,
			&code.ProfessionalID,
			&code.UsedAt,
			&code.ExpiresAt,
			&code.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		codes = append(codes, code)
	}

	return codes, rows.Err()
}

func (r *InvitationCodeRepository) GetInvitationCodeByID(
	ctx context.Context,
	id int,
) (models.InvitationCode, error) {
	query := `
		SELECT id, code, professional_id, used_at, expires_at, created_at
		FROM invitation_codes
		WHERE id = $1
	`

	var invCode models.InvitationCode
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&invCode.ID,
		&invCode.Code,
		&invCode.ProfessionalID,
		&invCode.UsedAt,
		&invCode.ExpiresAt,
		&invCode.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return models.InvitationCode{}, fmt.Errorf("invitation code not found")
	}
	if err != nil {
		return models.InvitationCode{}, err
	}

	return invCode, nil
}

func (r *InvitationCodeRepository) CreateInvitationCode(
	ctx context.Context,
	code string,
	professionalID int,
	expiresAt *time.Time,
) error {
	query := `
        INSERT INTO invitation_codes (code, professional_id, expires_at)
        VALUES ($1, $2, $3)
		`

	_, err := r.db.ExecContext(ctx, query, code, professionalID, expiresAt)
	return err
}

func (r *InvitationCodeRepository) InvalidateCode(
	ctx context.Context,
	id int,
) error {
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

func (r *InvitationCodeRepository) ReactivateCode(
	ctx context.Context,
	code string,
) error {
	query := `UPDATE invitation_codes SET used_at = NULL WHERE code = $1`

	result, err := r.db.ExecContext(ctx, query, code)
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
